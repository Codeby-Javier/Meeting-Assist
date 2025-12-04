package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Env                string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	JWTSecret          string
	JWTExpiration      int
	MaxAudioSize       int64
	MaxDocumentSize    int64
	MaxStoragePerUser  int64
	VoskModelEn        string
	TesseractPath      string
	FrontendURL        string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	AppConfig = &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "postgres"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		JWTExpiration:      getEnvAsInt("JWT_EXPIRATION", 24),
		MaxAudioSize:       getEnvAsInt64("MAX_AUDIO_SIZE", 104857600),
		MaxDocumentSize:    getEnvAsInt64("MAX_DOCUMENT_SIZE", 52428800),
		MaxStoragePerUser:  getEnvAsInt64("MAX_STORAGE_PER_USER", 1073741824),
		VoskModelEn:        getEnv("VOSK_MODEL_EN", "./vosk-model-en"),
		TesseractPath:      getEnv("TESSERACT_PATH", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		return value
	}
	return fallback
}
