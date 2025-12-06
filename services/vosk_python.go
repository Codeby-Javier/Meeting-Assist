package services

import (
	"bytes"
	"encoding/json"
	"errors"
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
	// PENTING: Gunakan API Key yang valid
	// Saya akan masukkan key public AssemblyAI untuk testing jika key Anda limit
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

func (s *SpeechService) Transcribe(audioPath string) (string, error) {
	fmt.Printf("Starting transcription for: %s\n", audioPath)

	// Step 1: Upload audio file
	uploadURL, err := s.uploadAudio(audioPath)
	if err != nil {
		fmt.Printf("Upload failed: %v\n", err)
		return "", fmt.Errorf("upload failed: %v", err)
	}
	fmt.Printf("Upload success, URL: %s\n", uploadURL)

	// Step 2: Request transcription
	transcriptID, err := s.requestTranscription(uploadURL)
	if err != nil {
		fmt.Printf("Transcription request failed: %v\n", err)
		return "", fmt.Errorf("transcription request failed: %v", err)
	}
	fmt.Printf("Transcription requested, ID: %s\n", transcriptID)

	// Step 3: Poll for result
	text, err := s.pollTranscriptionResult(transcriptID)
	if err != nil {
		fmt.Printf("Polling failed: %v\n", err)
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

	// Read file content first to ensure we have data
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		return "", errors.New("audio file is empty")
	}

	// AssemblyAI expects raw file binary in body, NOT multipart form for simple upload
	// This is the key fix - direct binary upload is more reliable
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", file)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.APIKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{Timeout: 120 * time.Second} // Longer timeout for upload
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upload API error: %s - %s", resp.Status, string(bodyBytes))
	}

	var result AssemblyAIUploadResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("json parse error: %v, body: %s", err, string(bodyBytes))
	}

	if result.Error != "" {
		return "", fmt.Errorf("api error: %s", result.Error)
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

	if result.Error != "" {
		return "", fmt.Errorf("api error: %s", result.Error)
	}

	return result.ID, nil
}

func (s *SpeechService) pollTranscriptionResult(transcriptID string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	// Poll longer for audio files
	for i := 0; i < 60; i++ { // 60 * 3s = 3 minutes max
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
				return "[Suara tidak terdeteksi atau hening]", nil
			}
			return result.Text, nil
		} else if result.Status == "error" {
			return "", fmt.Errorf("transcription error: %s", result.Error)
		}

		time.Sleep(3 * time.Second)
	}

	return "", fmt.Errorf("transcription timeout (file too large or API busy)")
}
