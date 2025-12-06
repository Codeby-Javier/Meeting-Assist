package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SpeechService struct {
	APIKey string
}

type AssemblyAIUploadResponse struct {
	UploadURL string `json:"upload_url"`
	Error     string `json:"error,omitempty"`
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
		APIKey: "4d2a2e6811d04ea99f0c5fa89090ddc6",
	}
}

func GetVoskService() *SpeechService {
	if speechService == nil {
		InitVosk()
	}
	return speechService
}

// Method 1: Upload langsung dari Stream (User -> Server -> AssemblyAI)
// Ini bypass penyimpanan disk lokal yang sering bermasalah
func (s *SpeechService) UploadStream(reader io.Reader) (string, error) {
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.APIKey)
	// Default binary upload

	client := &http.Client{Timeout: 0} // No timeout
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Network Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}

	var result AssemblyAIUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.UploadURL, nil
}

// Method 2: Request Transcript (Hanya butuh URL)
func (s *SpeechService) RequestTranscription(audioURL string) (string, error) {
	values := map[string]string{
		"audio_url":     audioURL,
		"language_code": "id",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AssemblyAITranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf(result.Error)
	}

	return result.ID, nil
}

// Method 3: Poll Result
func (s *SpeechService) PollResult(transcriptID string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	for i := 0; i < 200; i++ { // Poll 10 menit
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", s.APIKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		var result AssemblyAITranscriptResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.Status == "completed" {
			return result.Text, nil
		} else if result.Status == "error" {
			return "", fmt.Errorf(result.Error)
		}

		time.Sleep(3 * time.Second)
	}
	return "", fmt.Errorf("Timeout polling")
}

// Keep interface for compatibility if needed, but we will use new methods primarily
func (s *SpeechService) Transcribe(audioPath string) (string, error) {
	return "", fmt.Errorf("Deprecated: Use UploadStream -> Request -> Poll")
}
