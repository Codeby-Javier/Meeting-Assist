package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Environment
	if err := godotenv.Load(".env"); err != nil {
		// Try parent dir if running from scripts/
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Note: .env file not found, relying on system env vars")
		}
	}

	// 2. Connect DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	fmt.Println("Connecting to DB...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}

	fmt.Println("🚧 Sedang membersihkan database Supabase...")

	// 3. Hapus Data (Gunakan TRUNCATE CASCADE untuk menghapus relasi otomatis)
	tables := []string{
		"audio_transcripts",
		"document_texts",
		"audios",
		"documents",
		"users",
	}

	for _, table := range tables {
		// Cek apakah tabel exists
		if db.Migrator().HasTable(table) {
			err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error
			if err != nil {
				fmt.Printf("❌ Gagal membersihkan tabel %s: %v\n", table, err)
			} else {
				fmt.Printf("✅ Tabel %s BERSIH\n", table)
			}
		}
	}

	fmt.Println("✨ Selesai! Database sekarang kosong seperti baru.")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
