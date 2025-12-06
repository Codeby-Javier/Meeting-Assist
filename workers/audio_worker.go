package workers

import (
	"fmt"
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
)

func ProcessAudio(audioID uint, filePath string) {
	db := config.GetDB()
	speechService := services.GetVoskService()

	// Update status to processing
	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "processing")

	fmt.Printf("Processing audio ID %d, Path: %s\n", audioID, filePath)

	// Transcribe using AssemblyAI - No local conversion needed
	// AssemblyAI supports MP3, WAV, M4A, OGG etc natively
	text, err := speechService.Transcribe(filePath)

	if err != nil {
		fmt.Printf("Worker Error: %v\n", err)
		// Mark as failed
		db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "failed")

		// Optional: Save error message to text field for debug visibility
		// db.Model(&models.Audio{}).Where("id = ?", audioID).Update("filename", err.Error())
		return
	}

	// Save transcript
	transcript := models.AudioTranscript{
		AudioID: audioID,
		Text:    text,
	}

	if err := db.Create(&transcript).Error; err != nil {
		fmt.Printf("DB Save Error: %v\n", err)
		db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "failed")
		return
	}

	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "completed")
	fmt.Printf("Audio ID %d completed successfully\n", audioID)
}
