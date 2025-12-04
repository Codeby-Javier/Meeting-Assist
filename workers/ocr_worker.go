package workers

import (
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
	"meetingassist/utils"
	"os"
)

func ProcessDocument(docID uint, filePath string) {
	utils.LogDebug("=== Starting ProcessDocument for ID: %d, Path: %s ===", docID, filePath)
	db := config.GetDB()
	tesseractService := services.NewTesseractService()
	defer tesseractService.Close()

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.LogDebug("ERROR: File does not exist: %s", filePath)
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	// Update status to processing
	db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "processing")
	utils.LogDebug("Status updated to processing")

	// Preprocess image (this includes format conversion if needed)
	processedPath, err := utils.PreprocessImage(filePath)
	if err != nil {
		utils.LogDebug("PreprocessImage failed: %v, trying original file", err)
		// If preprocessing fails, try original
		processedPath = filePath
	} else {
		utils.LogDebug("Image preprocessed successfully: %s", processedPath)
	}

	// Verify processed file exists
	if _, err := os.Stat(processedPath); os.IsNotExist(err) {
		utils.LogDebug("ERROR: Processed file does not exist: %s", processedPath)
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	// OCR
	utils.LogDebug("Starting OCR on: %s", processedPath)
	text, err := tesseractService.PerformOCR(processedPath)
	if err != nil {
		// Log the specific error
		utils.LogDebug("OCR Failed for docID %d: %v", docID, err)
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	// Safe substring for logging
	logText := text
	if len(text) > 50 {
		logText = text[:50] + "..."
	}
	utils.LogDebug("OCR success! Text preview: %s", logText)

	// Save text
	docText := models.DocumentText{
		DocumentID: docID,
		Text:       text,
		PageNumber: 1,
	}
	if err := db.Create(&docText).Error; err != nil {
		utils.LogDebug("Failed to save document text to DB: %v", err)
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "completed")
	utils.LogDebug("=== ProcessDocument completed successfully for ID: %d ===", docID)
}
