package storage

import (
	"os"
)

// SaveToFile saves the given video data to the specified temporary file path.
func SaveToFile(tempFilePath string, videoData []byte) error {
	return os.WriteFile(tempFilePath, videoData, 0644)
}

// DeleteFile deletes the file at the specified temporary file path.
func DeleteFile(tempFilePath string) error {
	return os.Remove(tempFilePath)
}
