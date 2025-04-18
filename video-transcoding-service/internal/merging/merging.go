package merging

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Merger is responsible for merging video chunks.
type Merger struct {
	chunkDir       string
	outputFilePath string
	maxChunks      int
	checkInterval  time.Duration
	mu             sync.Mutex
	mergedChunks   map[int]bool
}

// NewMerger creates a new Merger instance.
func NewMerger(chunkDir string, outputFilePath string, maxChunks int, checkInterval time.Duration) *Merger {
	return &Merger{
		chunkDir:       chunkDir,
		outputFilePath: outputFilePath,
		maxChunks:      maxChunks,
		checkInterval:  checkInterval,
		mergedChunks:   make(map[int]bool),
	}
}

// Start begins the merging process.
func (m *Merger) Start() error {
	for {
		m.mu.Lock()
		err := m.mergeAvailableChunks()
		m.mu.Unlock()
		if err != nil {
			return err
		}

		if len(m.mergedChunks) == m.maxChunks {
			break
		}

		time.Sleep(m.checkInterval)
	}
	return nil
}

// mergeAvailableChunks merges all available chunks into the output file.
func (m *Merger) mergeAvailableChunks() error {
	outputFile, err := os.OpenFile(m.outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for i := 1; i <= m.maxChunks; i++ {
		if m.mergedChunks[i] {
			continue
		}

		chunkPath := filepath.Join(m.chunkDir, chunkFileName(i))
		if _, err := os.Stat(chunkPath); errors.Is(err, os.ErrNotExist) {
			continue
		}

		err := m.appendChunk(outputFile, chunkPath)
		if err != nil {
			return err
		}

		m.mergedChunks[i] = true
	}

	return nil
}

// appendChunk appends a single chunk to the output file.
func (m *Merger) appendChunk(outputFile *os.File, chunkPath string) error {
	chunkFile, err := os.Open(chunkPath)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	_, err = io.Copy(outputFile, chunkFile)
	return err
}

// chunkFileName generates the file name for a chunk based on its index.
func chunkFileName(index int) string {
	return "chunk" + strconv.Itoa(index)
}
