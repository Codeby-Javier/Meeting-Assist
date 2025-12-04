package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type StorageService struct {
	UploadDir string
}

func NewStorageService() *StorageService {
	return &StorageService{
		UploadDir: "uploads",
	}
}

func (s *StorageService) SaveFile(file *multipart.FileHeader, subDir string) (string, error) {
	// Create directory if not exists
	destDir := filepath.Join(s.UploadDir, subDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(destDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	return dst, nil
}
