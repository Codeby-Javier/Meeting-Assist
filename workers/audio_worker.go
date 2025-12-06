package workers

import (
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
	"meetingassist/utils"
)

func ProcessAudio(audioID uint, filePath string) {
	db := config.GetDB()
	voskService := services.GetVoskService()

	// Update status to processing
	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "processing")

	// Convert to WAV (no-op now, just returns original path)
	wavPath, _ := utils.ConvertToWav(filePath)

	// Transcribe using AssemblyAI
	text, err := voskService.Transcribe(wavPath)
	if err != nil {
		db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "failed")
		return
	}

	// Save transcript
	transcript := models.AudioTranscript{
		AudioID: audioID,
		Text:    text,
	}
	db.Create(&transcript)

	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "completed")
}
