package utils

// ConvertToWav is deprecated - cloud APIs accept multiple formats
// No need for local audio conversion anymore
func ConvertToWav(inputPath string) (string, error) {
	// Just return the original path - AssemblyAI accepts MP3, WAV, M4A, etc.
	return inputPath, nil
}
