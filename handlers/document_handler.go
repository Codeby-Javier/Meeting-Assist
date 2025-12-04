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

type DocumentHandler struct {
	Storage   *services.StorageService
	Tesseract *services.TesseractService
}

func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{
		Storage:   services.NewStorageService(),
		Tesseract: services.NewTesseractService(),
	}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No document uploaded")
		return
	}

	if file.Size > config.AppConfig.MaxDocumentSize {
		utils.ErrorResponse(c, http.StatusBadRequest, "File too large")
		return
	}

	path, err := h.Storage.SaveFile(file, "documents")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	userID := c.GetUint("user_id")
	doc := models.Document{
		UserID:   userID,
		Filename: file.Filename,
		FilePath: path,
		Type:     "image", // Simplified, should detect type
		Status:   "pending",
	}

	db := config.GetDB()
	if err := db.Create(&doc).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save record")
		return
	}

	// Trigger processing (async)
	go workers.ProcessDocument(doc.ID, path)

	utils.SuccessResponse(c, http.StatusAccepted, "Document uploaded and processing started", doc)
}

func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	userID := c.GetUint("user_id")
	var docs []models.Document
	config.GetDB().Where("user_id = ?", userID).Find(&docs)
	utils.SuccessResponse(c, http.StatusOK, "Documents retrieved", docs)
}

func (h *DocumentHandler) GetDocumentDetail(c *gin.Context) {
	id := c.Param("id")
	var doc models.Document
	if err := config.GetDB().Preload("Texts").First(&doc, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Document not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Document detail", doc)
}

func (h *DocumentHandler) UpdateText(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var docText models.DocumentText
	if err := config.GetDB().Where("document_id = ?", id).First(&docText).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Text not found")
		return
	}

	docText.Text = input.Text
	config.GetDB().Save(&docText)
	utils.SuccessResponse(c, http.StatusOK, "Text updated", docText)
}
