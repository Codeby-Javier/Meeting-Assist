package workers

import (
	"fmt"
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
)

// ProcessAudioURL handles transcription from an already uploaded Cloud URL
func ProcessAudioURL(audioID uint, audioURL string) {
	db := config.GetDB()
	speechService := services.GetVoskService()

	// Update status
	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "processing")

	// 1. Request Transcription
	transcriptID, err := speechService.RequestTranscription(audioURL)
	if err != nil {
		reportError(db, audioID, fmt.Sprintf("Request Transkrip Gagal: %v", err))
		return
	}

	// 2. Poll Result
	text, err := speechService.PollResult(transcriptID)
	if err != nil {
		reportError(db, audioID, fmt.Sprintf("Polling Gagal: %v", err))
		return
	}

	// Success
	transcript := models.AudioTranscript{
		AudioID: audioID,
		Text:    text,
	}
	db.Create(&transcript)

	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "completed")
}

// Deprecated: Old file-based processor (kept for interface compatibility if needed)
func ProcessAudio(audioID uint, filePath string) {
	// No-op
}

func reportError(db interface{}, audioID uint, msg string) {
	// Perlu type assertion atau akses DB langsung
	// Sederhananya kita asumsikan db adalah *gorm.DB dari context
	database := config.GetDB()

	database.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "failed")

	errorTranscript := models.AudioTranscript{
		AudioID: audioID,
		Text:    msg,
	}
	database.Create(&errorTranscript)
}
