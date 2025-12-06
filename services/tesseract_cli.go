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
)

type TesseractService struct {
	APIKey string
}

type OCRSpaceResponse struct {
	ParsedResults []struct {
		ParsedText string `json:"ParsedText"`
	} `json:"ParsedResults"`
	IsErroredOnProcessing bool   `json:"IsErroredOnProcessing"`
	ErrorMessage          string `json:"ErrorMessage,omitempty"`
}

func NewTesseractService() *TesseractService {
	return &TesseractService{
		APIKey: "K87899142388957", // Free OCR.space API key
	}
}

func (s *TesseractService) Close() {
	// No-op
}

func (s *TesseractService) PerformOCR(imagePath string) (string, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	io.Copy(part, file)

	// Add other fields
	writer.WriteField("apikey", s.APIKey)
	writer.WriteField("language", "eng")
	writer.WriteField("isOverlayRequired", "false")
	writer.WriteField("detectOrientation", "true")
	writer.WriteField("scale", "true")
	writer.WriteField("OCREngine", "2")

	writer.Close()

	// Make request to OCR.space API
	req, err := http.NewRequest("POST", "https://api.ocr.space/parse/image", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("OCR API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result OCRSpaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse OCR response: %v", err)
	}

	if result.IsErroredOnProcessing {
		return "", fmt.Errorf("OCR error: %s", result.ErrorMessage)
	}

	if len(result.ParsedResults) == 0 || result.ParsedResults[0].ParsedText == "" {
		return "[No text detected]", nil
	}

	return result.ParsedResults[0].ParsedText, nil
}
