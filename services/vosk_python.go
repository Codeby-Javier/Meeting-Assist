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
	speechService = &SpeechService{
		APIKey: "4d2a2e6811d04ea99f0c5fa89090ddc6", // Public Test Key
	}
}

func GetVoskService() *SpeechService {
	if speechService == nil {
		InitVosk()
	}
	return speechService
}

func (s *SpeechService) Transcribe(audioPath string) (string, error) {
	fmt.Printf("[ASSEMBLYAI] Starting processing for: %s\n", audioPath)

	// Cek file ada atau tidak
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return "", fmt.Errorf("File audio tidak ditemukan di server: %s", audioPath)
	}

	// 1. Upload
	uploadURL, err := s.uploadAudioChunked(audioPath)
	if err != nil {
		fmt.Printf("[ASSEMBLYAI] Upload Error: %v\n", err)
		return "", fmt.Errorf("Upload Gagal (Cek Koneksi): %v", err)
	}
	fmt.Printf("[ASSEMBLYAI] Upload Success: %s\n", uploadURL)

	// 2. Transcribe
	transcriptID, err := s.requestTranscription(uploadURL)
	if err != nil {
		fmt.Printf("[ASSEMBLYAI] Request Error: %v\n", err)
		return "", fmt.Errorf("Request API Gagal: %v", err)
	}
	fmt.Printf("[ASSEMBLYAI] Transcript ID: %s\n", transcriptID)

	// 3. Poll
	return s.pollTranscriptionResult(transcriptID)
}

// Upload menggunakan Chunked Transfer Encoding (Paling stabil untuk file besar)
func (s *SpeechService) uploadAudioChunked(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// AssemblyAI support standard binary upload.
	// Kuncinya adalah set header 'Transfer-Encoding: chunked' secara otomatis oleh Go
	// jika kita pass io.Reader (file) langsung ke http.NewRequest

	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", file)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.APIKey)
	// PENTING: Jangan set Content-Length manual, biarkan Go yang handle chunking
	// PENTING: Jangan set Content-Type multipart, pakai default binary

	client := &http.Client{
		Timeout: 0, // No timeout for upload (biarkan sampai selesai)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API Error %d: %s", resp.StatusCode, string(responseBytes))
	}

	var result AssemblyAIUploadResponse
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return "", fmt.Errorf("Parse Error: %v. Body: %s", err, string(responseBytes))
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Status %d: %s", resp.StatusCode, string(body))
	}

	var result AssemblyAITranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func (s *SpeechService) pollTranscriptionResult(transcriptID string) (string, error) {
	url := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	// Poll sampai 10 menit
	for i := 0; i < 200; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", s.APIKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue // Retry connection error
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
	return "", fmt.Errorf("Timeout: Proses terlalu lama (>10 menit)")
}
