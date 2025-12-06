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

	// Update status
	db.Model(&models.Audio{}).Where("id = ?", audioID).Update("status", "processing")

	// Transcribe
	text, err := speechService.Transcribe(filePath)

	if err != nil {
		fmt.Printf("AUDIO ERROR: %v\n", err)

		// PENTING: Simpan pesan error ke database agar terlihat di UI
		// Kita simpan object failed, dan text nya berisi pesan error
		db.Model(&models.Audio{}).Where("id = ?", audioID).Updates(map[string]interface{}{
			"status": "failed",
		})

		// Buat entry transcript yang berisi detail error
		// Jadi user bisa klik "Lihat" dan melihat errornya apa
		errorTranscript := models.AudioTranscript{
			AudioID: audioID,
			Text:    fmt.Sprintf("GAGAL: %v. Coba file yang lebih kecil atau format berbeda (MP3/WAV).", err),
		}
		db.Create(&errorTranscript)
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
