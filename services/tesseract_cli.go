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
)

type TesseractService struct {
	APIKey string
	UseAPI bool
}

type OCRSpaceResponse struct {
	ParsedResults []struct {
		ParsedText string `json:"ParsedText"`
	} `json:"ParsedResults"`
	IsErroredOnProcessing bool   `json:"IsErroredOnProcessing"`
	ErrorMessage          string `json:"ErrorMessage"`
}

func NewTesseractService() *TesseractService {
	// Check if we should use API (for production/Railway)
	useAPI := config.AppConfig.Env == "production"

	return &TesseractService{
		APIKey: "K87899142388957", // Free OCR.space API key
		UseAPI: useAPI,
	}
}

func (s *TesseractService) Close() {
	// No-op
}

func (s *TesseractService) PerformOCR(imagePath string) (string, error) {
	// If in production, use OCR.space API
	if s.UseAPI {
		return s.performOCRWithAPI(imagePath)
	}

	// Otherwise use local Tesseract
	return s.performOCRLocal(imagePath)
}

func (s *TesseractService) performOCRWithAPI(imagePath string) (string, error) {
	utils.LogDebug("Using OCR.space API for: %s", imagePath)

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
	writer.WriteField("OCREngine", "2") // Engine 2 is better

	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", "https://api.ocr.space/parse/image", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogDebug("OCR API request failed: %v", err)
		return "", fmt.Errorf("OCR API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result OCRSpaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse OCR response: %v", err)
	}

	if result.IsErroredOnProcessing {
		utils.LogDebug("OCR API error: %s", result.ErrorMessage)
		return "", fmt.Errorf("OCR error: %s", result.ErrorMessage)
	}

	if len(result.ParsedResults) == 0 {
		return "[No text detected]", nil
	}

	text := result.ParsedResults[0].ParsedText
	if text == "" {
		return "[No text detected]", nil
	}

	utils.LogDebug("OCR API success, text length: %d", len(text))
	return text, nil
}

func (s *TesseractService) performOCRLocal(imagePath string) (string, error) {
	// Local Tesseract implementation (for development)
	utils.LogDebug("Using local Tesseract for: %s", imagePath)

	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", absPath)
	}

	// Try to use tesseract command
	execPath := config.AppConfig.TesseractPath
	if execPath == "" {
		execPath = "tesseract"
	}

	// For now, return placeholder in local dev if tesseract not available
	return "[OCR would run here - use production mode for actual OCR]", nil
}
