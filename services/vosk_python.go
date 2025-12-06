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
	// Kunci Testing Public
	// Jika ini limit, User harus daftar di assemblyai.com (gratis)
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
	// Step 1: Upload (Multipart - Lebih Stabil)
	uploadURL, err := s.uploadAudioMultipart(audioPath)
	if err != nil {
		return "", fmt.Errorf("Upload Gagal: %v", err)
	}

	// Step 2: Request Transcription
	transcriptID, err := s.requestTranscription(uploadURL)
	if err != nil {
		return "", fmt.Errorf("Request Gagal: %v", err)
	}

	// Step 3: Poll
	return s.pollTranscriptionResult(transcriptID)
}

// Menggunakan Multipart Form Data (seperti OCR yang berhasil)
func (s *SpeechService) uploadAudioMultipart(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("file open error: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form file
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

	client := &http.Client{Timeout: 300 * time.Second} // 5 menit timeout upload
	resp, err := client.Do(req)
	if err != nil {
		// Cek koneksi internet server
		return "", fmt.Errorf("connection error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result AssemblyAIUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.UploadURL, nil
}

func (s *SpeechService) requestTranscription(audioURL string) (string, error) {
	values := map[string]string{
		"audio_url":     audioURL,
		"language_code": "id",
	}
	jsonData, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
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
		return "", fmt.Errorf(result.Error)
	}

	return result.ID, nil
}

func (s *SpeechService) pollTranscriptionResult(transcriptID string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	for i := 0; i < 60; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", s.APIKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
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
	return "", fmt.Errorf("Timeout waiting for transcription")
}
