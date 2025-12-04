package services

import (
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func (s *PDFService) CreatePDF(content string, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// Simple text wrapping
	pdf.MultiCell(0, 10, content, "", "", false)

	return pdf.OutputFileAndClose(outputPath)
}
