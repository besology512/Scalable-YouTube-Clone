package chunking

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// ChunkVideo chunks a video file into smaller parts if it exceeds 10 minutes or 1GB.
// It returns the path to the first chunk, the total number of chunks, and processes the rest concurrently.
// A channel is used to notify when a chunk is finished.
func ChunkVideo(inputFile string, outputDir string, chunkFinished chan<- string) (string, int, error) {
	// Ensure the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", 0, errors.New("input file does not exist")
	}

	// Get file size
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return "", 0, err
	}
	fileSizeGB := float64(fileInfo.Size()) / (1024 * 1024 * 1024)

	// Get video duration using ffprobe
	duration, err := getVideoDuration(inputFile)
	if err != nil {
		return "", 0, err
	}

	// Check if chunking is necessary
	if duration <= 600 && fileSizeGB <= 1 {
		return "", 0, errors.New("video does not need chunking")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", 0, err
	}

	// Determine the file extension of the input video
	fileExt := filepath.Ext(inputFile)
	if fileExt == "" {
		return "", 0, errors.New("unable to determine file extension")
	}

	// Calculate the total number of chunks
	totalChunks := int(duration / 600)
	if int(duration)%600 != 0 {
		totalChunks++
	}

	// Chunk the video using ffmpeg
	chunkPattern := filepath.Join(outputDir, "chunk_%03d"+fileExt)
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c", "copy", "-map", "0", "-segment_time", "600", "-f", "segment", chunkPattern)
	if err := cmd.Start(); err != nil {
		return "", 0, err
	}

	// Wait for the first chunk to be created
	firstChunk := filepath.Join(outputDir, "chunk_000"+fileExt)
	for {
		if _, err := os.Stat(firstChunk); err == nil {
			chunkFinished <- firstChunk // Notify that the first chunk is ready
			break
		}
	}

	// Process the rest of the chunks concurrently
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.Wait()

		// Notify for each chunk created
		for i := 1; i < totalChunks; i++ {
			chunkPath := filepath.Join(outputDir, "chunk_"+strings.Repeat("0", 3-len(strconv.Itoa(i)))+strconv.Itoa(i)+fileExt)
			if _, err := os.Stat(chunkPath); err == nil {
				chunkFinished <- chunkPath
			} else {
				break
			}
		}
		close(chunkFinished) // Close the channel when done
	}()

	return firstChunk, totalChunks, nil
}

// getVideoDuration uses ffprobe to get the duration of a video in seconds.
func getVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
