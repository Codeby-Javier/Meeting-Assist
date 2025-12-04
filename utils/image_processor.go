package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ConvertFormat converts HEIC/PDF/etc to PNG
func ConvertFormat(inputPath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		return inputPath, nil
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_converted.png"

	// For PDF, use Python script with pdf2image
	if ext == ".pdf" {
		absInput, _ := filepath.Abs(inputPath)
		absOutput, _ := filepath.Abs(outputPath)
		scriptPath, _ := filepath.Abs("scripts/pdf_to_png.py")

		LogDebug("Converting PDF: %s -> %s", absInput, absOutput)

		cmd := exec.Command("python", scriptPath, absInput, absOutput)
		output, err := cmd.CombinedOutput()

		if err != nil {
			LogDebug("PDF conversion failed: %v, output: %s", err, string(output))
			return "", fmt.Errorf("PDF conversion failed: %v", err)
		}

		LogDebug("PDF converted successfully: %s", absOutput)
		return outputPath, nil
	}

	// For HEIC and other formats, try ffmpeg
	absInput, _ := filepath.Abs(inputPath)
	absOutput, _ := filepath.Abs(outputPath)

	LogDebug("Converting format with ffmpeg: %s -> %s", absInput, absOutput)

	cmd := exec.Command("ffmpeg", "-i", absInput, "-y", absOutput)
	err := cmd.Run()
	if err != nil {
		LogDebug("Format conversion with ffmpeg failed: %v", err)
		return "", fmt.Errorf("conversion failed: %v", err)
	}

	LogDebug("Format converted successfully: %s", absOutput)
	return outputPath, nil
}

// PreprocessImage optimizes image for OCR with better handling for camera photos
func PreprocessImage(inputPath string) (string, error) {
	LogDebug("PreprocessImage called with: %s", inputPath)

	// First convert if necessary
	convertedPath, err := ConvertFormat(inputPath)
	if err != nil {
		LogDebug("ConvertFormat failed: %v, trying with original", err)
		convertedPath = inputPath
	} else {
		LogDebug("ConvertFormat succeeded: %s", convertedPath)
	}

	// Verify converted file exists
	if _, err := os.Stat(convertedPath); os.IsNotExist(err) {
		LogDebug("Converted file does not exist: %s", convertedPath)
		return "", fmt.Errorf("converted file does not exist: %s", convertedPath)
	}

	src, err := imaging.Open(convertedPath)
	if err != nil {
		LogDebug("imaging.Open failed: %v", err)
		return "", err
	}

	// Get image bounds to check size
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	LogDebug("Original image size: %dx%d", width, height)

	// Resize if image is too large (better for OCR and performance)
	maxDimension := 2000
	if width > maxDimension || height > maxDimension {
		if width > height {
			src = imaging.Resize(src, maxDimension, 0, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, 0, maxDimension, imaging.Lanczos)
		}
		LogDebug("Resized image to fit max dimension: %d", maxDimension)
	}

	// Apply OCR-friendly preprocessing
	// 1. Grayscale
	src = imaging.Grayscale(src)

	// 2. Increase contrast (helps with camera photos)
	src = imaging.AdjustContrast(src, 30)

	// 3. Adjust brightness slightly
	src = imaging.AdjustBrightness(src, 5)

	// 4. Sharpen (helps with slightly blurry photos)
	src = imaging.Sharpen(src, 1.5)

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_processed.png"
	absOutput, _ := filepath.Abs(outputPath)

	// Save with high quality
	err = imaging.Save(src, absOutput, imaging.PNGCompressionLevel(png.BestCompression))
	if err != nil {
		LogDebug("imaging.Save failed: %v", err)
		return "", err
	}

	LogDebug("PreprocessImage completed: %s", absOutput)
	return absOutput, nil
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
