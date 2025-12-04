package handlers

import (
	"meetingassist/services"
	"meetingassist/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	PDF  *services.PDFService
	Docx *services.DocxService
}

func NewExportHandler() *ExportHandler {
	return &ExportHandler{
		PDF:  &services.PDFService{},
		Docx: &services.DocxService{},
	}
}

func (h *ExportHandler) ExportPDF(c *gin.Context) {
	var input struct {
		Content  string `json:"content" binding:"required"`
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	outputPath := filepath.Join("uploads", "exports", input.Filename+".pdf")
	if err := h.PDF.CreatePDF(input.Content, outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create PDF")
		return
	}

	c.File(outputPath)
}

func (h *ExportHandler) ExportDocx(c *gin.Context) {
	var input struct {
		Content  string `json:"content" binding:"required"`
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	outputPath := filepath.Join("uploads", "exports", input.Filename+".docx")
	if err := h.Docx.CreateDocx(input.Content, outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create DOCX")
		return
	}

	c.File(outputPath)
}
