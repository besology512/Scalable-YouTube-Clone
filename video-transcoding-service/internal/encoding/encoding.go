package encoding

import (
	"fmt"
	"os/exec"
	"time"
)

const maxRetries = 3
const retryDelay = 2 * time.Second

// TranscodeToHLS transcodes a video to HLS format.
func TranscodeToHLS(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-codec: copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", outputPath)
		return cmd.Run()
	}, "transcode to HLS")
}

// TranscodeTo144p transcodes a video to MP4 format with 144p quality.
func TranscodeTo144p(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=256:144", "-c:v", "libx264", "-preset", "slow", "-crf", "28", "-c:a", "aac", "-b:a", "64k", outputPath)
		return cmd.Run()
	}, "transcode to 144p")
}

// TranscodeTo360p transcodes a video to MP4 format with 360p quality.
func TranscodeTo360p(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=640:360", "-c:v", "libx264", "-preset", "slow", "-crf", "23", "-c:a", "aac", "-b:a", "128k", outputPath)
		return cmd.Run()
	}, "transcode to 360p")
}

// TranscodeTo720p transcodes a video to MP4 format with 720p quality.
func TranscodeTo720p(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=1280:720", "-c:v", "libx264", "-preset", "slow", "-crf", "20", "-c:a", "aac", "-b:a", "192k", outputPath)
		return cmd.Run()
	}, "transcode to 720p")
}

// TranscodeTo1080p transcodes a video to MP4 format with 1080p quality.
func TranscodeTo1080p(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=1920:1080", "-c:v", "libx264", "-preset", "slow", "-crf", "18", "-c:a", "aac", "-b:a", "256k", outputPath)
		return cmd.Run()
	}, "transcode to 1080p")
}

// TranscodeTo4K transcodes a video to MP4 format with 4K quality.
func TranscodeTo4K(inputPath, outputPath string) error {
	return retry(func() error {
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=3840:2160", "-c:v", "libx264", "-preset", "slow", "-crf", "15", "-c:a", "aac", "-b:a", "320k", outputPath)
		return cmd.Run()
	}, "transcode to 4K")
}

// retry retries a function up to maxRetries times with a delay between attempts.
func retry(action func() error, actionName string) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = action(); err == nil {
			return nil
		}
		fmt.Printf("Attempt %d to %s failed: %v. Retrying...\n", i+1, actionName, err)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("all attempts to %s failed: %w", actionName, err)
}
