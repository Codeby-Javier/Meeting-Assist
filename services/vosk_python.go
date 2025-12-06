package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	// Ambil API Key dari Environment Variable
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	// Fallback ke key public (sering error/limit)
	if apiKey == "" {
		// Log warning di server
		fmt.Println("[WARNING] ASSEMBLYAI_API_KEY tidak ditemukan di environment variables. Menggunakan shared key (mungkin expired).")
		apiKey = "4d2a2e6811d04ea99f0c5fa89090ddc6"
	} else {
		fmt.Println("[INFO] Menggunakan Custom ASSEMBLYAI_API_KEY")
	}

	speechService = &SpeechService{
		APIKey: apiKey,
	}
}

func GetVoskService() *SpeechService {
	if speechService == nil {
		InitVosk()
	}
	return speechService
}

// Method 1: Upload Langsung (Chunked)
func (s *SpeechService) UploadStream(reader io.Reader) (string, error) {
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.APIKey)

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Network Error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("API Key Invalid/Expired. Harap set 'ASSEMBLYAI_API_KEY' di Railway Variables.")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}

	var result AssemblyAIUploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.UploadURL, nil
}

// Method 2: Request Transcript
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

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("API Key Invalid/Expired. Harap set 'ASSEMBLYAI_API_KEY' di Railway Variables.")
	}

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

	for i := 0; i < 200; i++ {
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

// Interface compatibility
func (s *SpeechService) Transcribe(audioPath string) (string, error) {
	return "", fmt.Errorf("Use UploadStream instead")
}
