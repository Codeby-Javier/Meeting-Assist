package services

import (
	"bytes"
	"fmt"
	"meetingassist/config"
	"meetingassist/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TesseractService struct {
	ExecutablePath string
}

func NewTesseractService() *TesseractService {
	path := config.AppConfig.TesseractPath
	if path == "" {
		path = "tesseract" // Default to PATH
	}
	return &TesseractService{
		ExecutablePath: path,
	}
}

func (s *TesseractService) Close() {
	// No-op for CLI
}

func (s *TesseractService) PerformOCR(imagePath string) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		utils.LogDebug("Failed to get absolute path for %s: %v", imagePath, err)
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		utils.LogDebug("File does not exist: %s", absPath)
		return "", fmt.Errorf("file does not exist: %s", absPath)
	}

	utils.LogDebug("Running Tesseract on: %s", absPath)

	// tesseract <image> stdout
	cmd := exec.Command(s.ExecutablePath, absPath, "stdout")

	// Capture output
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		utils.LogDebug("Tesseract failed: %v, stderr: %s", err, stderrStr)
		return "", fmt.Errorf("tesseract error: %v, stderr: %s", err, stderrStr)
	}

	text := strings.TrimSpace(out.String())
	utils.LogDebug("Tesseract output length: %d chars", len(text))

	if text == "" {
		return "[No text detected]", nil
	}

	return text, nil
}
