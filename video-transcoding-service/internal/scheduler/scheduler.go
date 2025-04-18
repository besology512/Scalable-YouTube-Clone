package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"video-transcoding-service/internal/chunking"
	"video-transcoding-service/internal/config"
	"video-transcoding-service/internal/encoding"
	"video-transcoding-service/internal/kafka"
	"video-transcoding-service/internal/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Scheduler struct {
	cfg         *config.Config
	kafkaProd   *kafka.Producer
	minioClient *minio.Client
	consumer    *kafka.Consumer // Add consumer to struct
}

// Add at the top of the file
const maxConcurrentProcesses = 5

func New(cfg *config.Config) (*Scheduler, error) {
	// 1. Initialize Kafka producer first
	prod, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 2. Initialize MinIO client
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// 3. Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &Scheduler{
		cfg:         cfg,
		kafkaProd:   prod,
		minioClient: minioClient,
		consumer:    consumer,
	}, nil
}

// Update the Run method
func (s *Scheduler) Run() error {
	sem := make(chan struct{}, maxConcurrentProcesses)
	defer s.consumer.Close()

	ctx := context.Background()
	go func() {
		if err := s.consumer.Consume(ctx, []string{s.cfg.KafkaTranscodeTopic}); err != nil {
			log.Printf("Error starting consumer: %v", err)
		}
	}()

	for msg := range s.consumer.Messages() {
		sem <- struct{}{} // Acquire slot
		go func(m *types.VideoMessage) {
			defer func() { <-sem }()
			s.processVideoMessage(m)
		}(msg)
	}
	return nil
}

// Rename to processVideoMessage and change parameter
func (s *Scheduler) processVideoMessage(msg *types.VideoMessage) {
	// Add validation for the message
	if msg == nil || msg.VideoURL == "" {
		log.Println("Received invalid video message")
		return
	}

	// Extract object name from URL
	parts := strings.Split(msg.VideoURL, "/")
	if len(parts) == 0 {
		log.Printf("Invalid video URL format: %s", msg.VideoURL)
		return
	}
	objectName := parts[len(parts)-1]
	localPath, cleanup, err := s.downloadVideo(objectName)
	if err != nil {
		log.Printf("Failed to download video: %v", err)
		return
	}
	defer cleanup()

	workspace, err := os.MkdirTemp("", "video-processing-")
	if err != nil {
		log.Printf("Failed to create workspace: %v", err)
		return
	}
	defer os.RemoveAll(workspace)

	originalFileName := filepath.Base(objectName)
	baseName := strings.TrimSuffix(originalFileName, filepath.Ext(originalFileName))

	chunks, totalChunks, err := s.chunkVideoIfNeeded(localPath, workspace)
	if err != nil {
		log.Printf("Chunking failed: %v", err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan string)

	resolutions := []struct {
		name   string
		method func(string, string) error
	}{
		{"144p", encoding.TranscodeTo144p},
		{"360p", encoding.TranscodeTo360p},
		{"720p", encoding.TranscodeTo720p},
		{"1080p", encoding.TranscodeTo1080p},
		{"4K", encoding.TranscodeTo4K},
	}

	for _, res := range resolutions {
		wg.Add(1)
		go func(resName string, transcodeFunc func(string, string) error) {
			defer wg.Done()
			s.processResolution(chunks, totalChunks, workspace, resName, transcodeFunc, results)
		}(res.name, res.method)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processHLS(localPath, workspace, baseName, results)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		relPath, err := filepath.Rel(workspace, result)
		if err != nil {
			log.Printf("Failed to get relative path for %s: %v", result, err)
			continue
		}
		objectKey := filepath.Join(baseName, relPath)
		if err := s.uploadToMinio(result, objectKey); err != nil {
			log.Printf("Failed to upload %s: %v", result, err)
			continue
		}
		s.kafkaProd.ProduceMessage(s.cfg.KafkaUploadTranscodedTopic, objectKey)
	}
}

func (s *Scheduler) downloadVideo(objectName string) (string, func(), error) {
	// Generate a temporary file path WITHOUT creating the file
	tmpDir := os.TempDir()
	tmpPattern := fmt.Sprintf("video-%s-*.mp4", objectName)
	tmpFile, err := os.CreateTemp(tmpDir, tmpPattern)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Immediately close the file handle - MinIO will manage the file
	tmpFile.Close()

	// Download directly to the temp path
	err = s.minioClient.FGetObject(context.Background(),
		s.cfg.MinioRawVideosBucket,
		objectName,
		tmpPath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("failed to download from MinIO: %w", err)
	}

	cleanup := func() {
		if err := os.Remove(tmpPath); err != nil {
			log.Printf("Failed to cleanup temp file %s: %v", tmpPath, err)
		}
	}

	return tmpPath, cleanup, nil
}

func (s *Scheduler) uploadToMinio(filePath, objectKey string) error {
	_, err := s.minioClient.FPutObject(context.Background(), s.cfg.MinioProcessedVideosBucket, objectKey, filePath, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload to MinIO: %w", err)
	}
	return nil
}

func (s *Scheduler) processResolution(chunks []string, totalChunks int, workspace, resName string,
	transcodeFunc func(string, string) error, results chan<- string) {

	resDir := filepath.Join(workspace, resName)
	if err := os.MkdirAll(resDir, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v", resDir, err)
		return
	}

	var transcodedChunks []string
	for i, chunk := range chunks {
		outputPath := filepath.Join(resDir, fmt.Sprintf("chunk_%d.mp4", i))
		if err := transcodeFunc(chunk, outputPath); err != nil {
			log.Printf("Failed to transcode %s chunk %d: %v", resName, i, err)
			return
		}
		transcodedChunks = append(transcodedChunks, outputPath)
	}

	mergedPath := filepath.Join(workspace, fmt.Sprintf("%s-final.mp4", resName))
	if err := s.mergeChunks(transcodedChunks, mergedPath); err != nil {
		log.Printf("Failed to merge %s chunks: %v", resName, err)
		return
	}

	results <- mergedPath
}

func (s *Scheduler) processHLS(inputPath, workspace, baseName string, results chan<- string) {
	hlsDir := filepath.Join(workspace, "hls")
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		log.Printf("Failed to create HLS directory: %v", err)
		return
	}

	hlsPath := filepath.Join(hlsDir, "output.m3u8")
	if err := encoding.TranscodeToHLS(inputPath, hlsPath); err != nil {
		log.Printf("HLS transcoding failed: %v", err)
		return
	}

	filepath.Walk(hlsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			results <- path
		}
		return nil
	})
}

func (s *Scheduler) chunkVideoIfNeeded(inputPath, outputDir string) ([]string, int, error) {
	chunkFinished := make(chan string, 10)
	firstChunk, totalChunks, err := chunking.ChunkVideo(inputPath, outputDir, chunkFinished)
	if err != nil && err.Error() != "video does not need chunking" {
		return nil, 0, err
	}

	if totalChunks == 0 {
		return []string{inputPath}, 1, nil
	}

	chunks := make([]string, totalChunks)
	chunks[0] = firstChunk

	timeout := time.After(30 * time.Minute)
	for i := 1; i < totalChunks; i++ {
		select {
		case chunk := <-chunkFinished:
			chunks[i] = chunk
		case <-timeout:
			return nil, 0, fmt.Errorf("timed out waiting for chunks")
		}
	}

	return chunks, totalChunks, nil
}

func (s *Scheduler) mergeChunks(chunks []string, outputPath string) error {
	listFile := outputPath + ".txt"
	defer os.Remove(listFile)

	f, err := os.Create(listFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, chunk := range chunks {
		_, err := f.WriteString(fmt.Sprintf("file '%s'\n", chunk))
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", listFile, "-c", "copy", outputPath)
	return cmd.Run()
}
