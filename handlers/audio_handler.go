package handlers

import (
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/services"
	"meetingassist/utils"
	"meetingassist/workers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AudioHandler struct {
	Storage *services.StorageService
	Speech  *services.SpeechService
}

func NewAudioHandler() *AudioHandler {
	return &AudioHandler{
		Storage: services.NewStorageService(),
		Speech:  services.GetVoskService(),
	}
}

func (h *AudioHandler) UploadAudio(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No audio file uploaded")
		return
	}

	if file.Size > config.AppConfig.MaxAudioSize {
		utils.ErrorResponse(c, http.StatusBadRequest, "File too large")
		return
	}

	path, err := h.Storage.SaveFile(file, "audio")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	userID := c.GetUint("user_id")
	audio := models.Audio{
		UserID:   userID,
		Filename: file.Filename,
		FilePath: path,
		Status:   "pending",
	}

	db := config.GetDB()
	if err := db.Create(&audio).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save record")
		return
	}

	// Trigger processing (async)
	go workers.ProcessAudio(audio.ID, path)

	utils.SuccessResponse(c, http.StatusAccepted, "Audio uploaded and processing started", audio)
}

func (h *AudioHandler) GetAudios(c *gin.Context) {
	userID := c.GetUint("user_id")
	var audios []models.Audio
	config.GetDB().Where("user_id = ?", userID).Find(&audios)
	utils.SuccessResponse(c, http.StatusOK, "Audios retrieved", audios)
}

func (h *AudioHandler) GetAudioDetail(c *gin.Context) {
	id := c.Param("id")
	var audio models.Audio
	if err := config.GetDB().Preload("Transcripts").First(&audio, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Audio not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Audio detail", audio)
}

func (h *AudioHandler) UpdateTranscript(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update the first transcript for simplicity
	// In a real app, we might have multiple segments
	var transcript models.AudioTranscript
	if err := config.GetDB().Where("audio_id = ?", id).First(&transcript).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transcript not found")
		return
	}

	transcript.Text = input.Text
	config.GetDB().Save(&transcript)
	utils.SuccessResponse(c, http.StatusOK, "Transcript updated", transcript)
}
