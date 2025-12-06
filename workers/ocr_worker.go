package workers

import (
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
	"meetingassist/utils"
	"os"
)

func ProcessDocument(docID uint, filePath string) {
	db := config.GetDB()
	tesseractService := services.NewTesseractService()
	defer tesseractService.Close()

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	// Update status to processing
	db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "processing")

	// Preprocess image (minimal - just resize if needed)
	processedPath, err := utils.PreprocessImage(filePath)
	if err != nil {
		processedPath = filePath // Use original if preprocess fails
	}

	// OCR using OCR.space API
	text, err := tesseractService.PerformOCR(processedPath)
	if err != nil {
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	// Save text
	docText := models.DocumentText{
		DocumentID: docID,
		Text:       text,
		PageNumber: 1,
	}
	if err := db.Create(&docText).Error; err != nil {
		db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "failed")
		return
	}

	db.Model(&models.Document{}).Where("id = ?", docID).Update("status", "completed")
}
