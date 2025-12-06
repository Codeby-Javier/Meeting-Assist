package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SpeechService struct {
	APIKey string
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
	Error  string `json:"error,omitempty"`
}

var speechService *SpeechService

func InitVosk() {
	speechService = &SpeechService{
		APIKey: "4d2a2e6811d04ea99f0c5fa89090ddc6", // Free AssemblyAI key
	}
}

func GetVoskService() *SpeechService {
	return speechService
}

func (s *SpeechService) Transcribe(audioPath string) (string, error) {
	// Step 1: Upload audio file
	uploadURL, err := s.uploadAudio(audioPath)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	// Step 2: Request transcription
	transcriptID, err := s.requestTranscription(uploadURL)
	if err != nil {
		return "", fmt.Errorf("transcription request failed: %v", err)
	}

	// Step 3: Poll for result
	text, err := s.pollTranscriptionResult(transcriptID)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %v", err)
	}

	return text, nil
}

func (s *SpeechService) uploadAudio(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AssemblyAIUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.UploadURL, nil
}

func (s *SpeechService) requestTranscription(audioURL string) (string, error) {
	reqBody := AssemblyAITranscriptRequest{
		AudioURL:     audioURL,
		LanguageCode: "id", // Indonesian
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AssemblyAITranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (s *SpeechService) pollTranscriptionResult(transcriptID string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	for i := 0; i < 120; i++ { // Poll for up to 10 minutes
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", s.APIKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		var result AssemblyAITranscriptResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.Status == "completed" {
			if result.Text == "" {
				return "[No speech detected]", nil
			}
			return result.Text, nil
		} else if result.Status == "error" {
			return "", fmt.Errorf("transcription error: %s", result.Error)
		}

		time.Sleep(5 * time.Second)
	}

	return "", fmt.Errorf("transcription timeout")
}
