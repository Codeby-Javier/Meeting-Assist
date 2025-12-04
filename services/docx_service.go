package services

import (
	"os"
)

type DocxService struct{}

func (s *DocxService) CreateDocx(content string, outputPath string) error {
	// Creating a real DOCX is complex without a heavy library.
	// We will write an HTML file which Word can open, but save it with .doc extension (or .docx)
	// This is a common workaround for simple requirements.

	htmlContent := `
	<html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word' xmlns='http://www.w3.org/TR/REC-html40'>
	<head><meta charset='utf-8'><title>Export</title></head>
	<body>
	` + content + `
	</body></html>
	`

	return os.WriteFile(outputPath, []byte(htmlContent), 0644)
}
