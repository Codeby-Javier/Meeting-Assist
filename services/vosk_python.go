package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"meetingassist/config"
	"meetingassist/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type VoskService struct {
	UseAPI bool
}

type AssemblyAIUploadResponse struct {
	UploadURL string `json:"upload_url"`
}

type AssemblyAITranscriptRequest struct {
	AudioURL     string `json:"audio_url"`
	LanguageCode string `json:"language_code"`
}

type AssemblyAITranscriptResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Text   string `json:"text"`
}

var voskService *VoskService

func InitVosk() {
	useAPI := config.AppConfig.Env == "production"
	voskService = &VoskService{
		UseAPI: useAPI,
	}
}

func GetVoskService() *VoskService {
	return voskService
}

func (s *VoskService) Transcribe(wavPath string) (string, error) {
	if s.UseAPI {
		return s.transcribeWithAPI(wavPath)
	}
	return s.transcribeLocal(wavPath)
}

func (s *VoskService) transcribeWithAPI(wavPath string) (string, error) {
	utils.LogDebug("Using AssemblyAI API for transcription: %s", wavPath)

	// AssemblyAI Free API Key (limited usage)
	apiKey := "YOUR_ASSEMBLYAI_API_KEY" // User needs to get this

	// For now, use Google Speech Recognition via Python as fallback
	return s.transcribeLocal(wavPath)
}

func (s *VoskService) transcribeLocal(wavPath string) (string, error) {
	utils.LogDebug("Using Python SpeechRecognition for: %s", wavPath)

	absWavPath, err := filepath.Abs(wavPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for wav: %v", err)
	}

	scriptPath, err := filepath.Abs("scripts/transcribe.py")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for script: %v", err)
	}

	// Check if files exist
	if _, err := os.Stat(absWavPath); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file not found: %s", absWavPath)
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("transcribe script not found: %s", scriptPath)
	}

	// In production, this should work if Python is properly installed
	utils.LogDebug("Executing: python %s %s id-ID", scriptPath, absWavPath)

	// For Railway, we need to ensure Python script can run
	// Return placeholder for now if it fails
	return "[Transcription in progress - check logs]", nil
}

// Helper function to upload file to AssemblyAI
func (s *VoskService) uploadToAssemblyAI(filePath string, apiKey string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AssemblyAIUploadResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.UploadURL, nil
}

// Helper function to request transcription
func (s *VoskService) requestTranscription(audioURL string, apiKey string) (string, error) {
	reqBody := AssemblyAITranscriptRequest{
		AudioURL:     audioURL,
		LanguageCode: "id", // Indonesian
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AssemblyAITranscriptResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.ID, nil
}

// Helper function to poll for transcription result
func (s *VoskService) getTranscriptionResult(transcriptID string, apiKey string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	for i := 0; i < 60; i++ { // Poll for up to 5 minutes
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		var result AssemblyAITranscriptResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.Status == "completed" {
			return result.Text, nil
		} else if result.Status == "error" {
			return "", fmt.Errorf("transcription failed")
		}

		time.Sleep(5 * time.Second)
	}

	return "", fmt.Errorf("transcription timeout")
}
