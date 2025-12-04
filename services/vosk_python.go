package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type VoskService struct{}

var voskService *VoskService

func InitVosk() {
	voskService = &VoskService{}
}

func GetVoskService() *VoskService {
	return voskService
}

func (s *VoskService) Transcribe(wavPath string) (string, error) {
	absWavPath, err := filepath.Abs(wavPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for wav: %v", err)
	}

	scriptPath, err := filepath.Abs("scripts/transcribe.py")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for script: %v", err)
	}

	// Call python script
	// python <scriptPath> <absWavPath> id-ID
	cmd := exec.Command("python", scriptPath, absWavPath, "id-ID")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("transcription error: %v, stderr: %s", err, stderr.String())
	}

	text := strings.TrimSpace(out.String())
	return text, nil
}
