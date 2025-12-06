package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// PreprocessImage optimizes image for OCR (minimal processing)
func PreprocessImage(inputPath string) (string, error) {
	// Just ensure it's a supported format and resize if too large
	ext := strings.ToLower(filepath.Ext(inputPath))

	// If PNG or JPG, just resize if needed
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		src, err := imaging.Open(inputPath)
		if err != nil {
			return inputPath, nil // Return original if can't open
		}

		bounds := src.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Resize if too large (OCR.space has 1MB limit)
		maxDimension := 2000
		if width > maxDimension || height > maxDimension {
			if width > height {
				src = imaging.Resize(src, maxDimension, 0, imaging.Lanczos)
			} else {
				src = imaging.Resize(src, 0, maxDimension, imaging.Lanczos)
			}

			outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_resized.png"
			err = imaging.Save(src, outputPath)
			if err != nil {
				return inputPath, nil
			}
			return outputPath, nil
		}

		return inputPath, nil
	}

	// For other formats, try to convert to PNG
	src, err := imaging.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("unsupported image format: %v", err)
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_converted.png"
	err = imaging.Save(src, outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func SaveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".jpg" || ext == ".jpeg" {
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	}
	return png.Encode(f, img)
}
