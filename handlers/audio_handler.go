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
	Speech *services.SpeechService
}

func NewAudioHandler() *AudioHandler {
	return &AudioHandler{
		Speech: services.GetVoskService(),
	}
}

func (h *AudioHandler) UploadAudio(c *gin.Context) {
	// 1. Terima File
	fileHeader, err := c.FormFile("audio")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No audio file uploaded")
		return
	}

	// 2. Upload Langsung ke AssemblyAI (Tanpa Simpan Lokal)
	srcFile, err := fileHeader.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to open file stream")
		return
	}
	defer srcFile.Close()

	// Upload Stream
	uploadURL, err := h.Speech.UploadStream(srcFile)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Upload ke Cloud Gagal: "+err.Error())
		return
	}

	// 3. Simpan Record ke DB
	userID := c.GetUint("user_id")
	audio := models.Audio{
		UserID:   userID,
		Filename: fileHeader.Filename,
		FilePath: uploadURL, // Simpan URL Cloud, bukan path lokal!
		Status:   "uploaded",
	}

	db := config.GetDB()
	if err := db.Create(&audio).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save record")
		return
	}

	// 4. Trigger Worker untuk Polling (Kirim URL Cloud)
	go workers.ProcessAudioURL(audio.ID, uploadURL)

	utils.SuccessResponse(c, http.StatusAccepted, "Audio uploaded to cloud, processing started", audio)
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

	var transcript models.AudioTranscript
	if err := config.GetDB().Where("audio_id = ?", id).First(&transcript).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transcript not found")
		return
	}

	transcript.Text = input.Text
	config.GetDB().Save(&transcript)
	utils.SuccessResponse(c, http.StatusOK, "Transcript updated", transcript)
}
