package ffmpeg

import (
	"log"
	"os/exec"
)

// Transcode converts the uploaded video into a standard format
func Transcode(inputPath string, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx264", "-acodec", "aac", outputPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error transcoding video: %v", err)
		return err
	}
	return nil
}
