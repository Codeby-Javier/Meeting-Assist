package utils

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertToWav converts audio to WAV format (16kHz, Mono) required by Vosk
func ConvertToWav(inputPath string) (string, error) {
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".wav"

	absInput, _ := filepath.Abs(inputPath)
	absOutput, _ := filepath.Abs(outputPath)

	// ffmpeg -i input.mp3 -ar 16000 -ac 1 output.wav
	cmd := exec.Command("ffmpeg", "-i", absInput, "-ar", "16000", "-ac", "1", "-y", absOutput)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
