package workers

import (
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
	"meetingassist/utils"
	"path/filepath"
)

func ProcessAudio(audioID uint, filePath string) {
	utils.LogDebug("Starting ProcessAudio for ID: %d, Path: %s", audioID, filePath)
	db := config.GetDB()
	voskService := services.GetVoskService()

	// Update status to processing
	if err := db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "processing").Error; err != nil {
		utils.LogDebug("Failed to update status to processing: %v", err)
	}

	// Convert to WAV
	wavPath, err := utils.ConvertToWav(filePath)
	if err != nil {
		utils.LogDebug("ConvertToWav failed: %v", err)
		// If conversion fails, try original if it's wav, else fail
		if filepath.Ext(filePath) != ".wav" {
			// Try to proceed with original path if conversion fails, maybe it works?
			// But let's log it.
		}
		wavPath = filePath // Fallback
	}
	utils.LogDebug("WavPath: %s", wavPath)

	// Transcribe
	text, err := voskService.Transcribe(wavPath)
	if err != nil {
		utils.LogDebug("Transcription Failed for audioID %d: %v", audioID, err)
		db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "failed")
		return
	}

	// Safe substring for logging
	logText := text
	if len(text) > 20 {
		logText = text[:20] + "..."
	}
	utils.LogDebug("Transcription success: %s", logText)

	// Save transcript
	transcript := models.AudioTranscript{
		AudioID: audioID,
		Text:    text,
	}
	if err := db.Create(&transcript).Error; err != nil {
		utils.LogDebug("Failed to save transcript to DB: %v", err)
	}

	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "completed")
	utils.LogDebug("ProcessAudio completed for ID: %d", audioID)
}
