package quality_check

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// DetectQuality takes a video file path and returns its quality (e.g., 720, 1080, etc.)
func DetectQuality(filePath string) (int, error) {
	// Use ffprobe to extract video stream information
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=height", "-of", "csv=p=0", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, errors.New("failed to run ffprobe: " + err.Error())
	}

	// Parse the height from the output
	heightStr := strings.TrimSpace(string(output))
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return 0, errors.New("failed to parse video height: " + err.Error())
	}

	return height, nil
}

// IsValidQuality checks if the detected quality is one of the standard qualities
func IsValidQuality(quality int) bool {
	validQualities := []int{144, 240, 360, 480, 720, 1080, 1440, 2160}
	for _, q := range validQualities {
		if quality == q {
			return true
		}
	}
	return false
}

// Example usage
func Example() {
	filePath := "example_video.mp4"
	quality, err := DetectQuality(filePath)
	if err != nil {
		panic("Error detecting quality: " + err.Error())
	}

	if IsValidQuality(quality) {
		println("Detected quality:", quality)
	} else {
		println("Detected quality is non-standard:", quality)
	}
}
