package config

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var dialector gorm.Dialector

	if AppConfig.DBHost == "sqlite" {
		dialector = sqlite.Open("meetingassist.db")
		log.Println("Using SQLite database")
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			AppConfig.DBHost,
			AppConfig.DBUser,
			AppConfig.DBPassword,
			AppConfig.DBName,
			AppConfig.DBPort,
			AppConfig.DBSSLMode,
		)
		dialector = postgres.Open(dsn)
		log.Println("Using PostgreSQL database")
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println("Failed to connect to configured database:", err)
		log.Println("Falling back to local SQLite database (meetingassist.db)...")
		DB, err = gorm.Open(sqlite.Open("meetingassist.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Failed to connect to database: ", err)
		}
	}

	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}
