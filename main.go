package main

import (
	"log"
	"meetingassist/config"
	"meetingassist/middleware"
	"meetingassist/models"
	"meetingassist/routes"
	"meetingassist/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Connect Database
	config.ConnectDB()

	// Ensure directories exist
	os.MkdirAll("uploads/audio", 0755)
	os.MkdirAll("uploads/documents", 0755)
	os.MkdirAll("uploads/exports", 0755)

	// Migrate Models
	db := config.GetDB()
	err := db.AutoMigrate(&models.User{}, &models.Audio{}, &models.AudioTranscript{}, &models.Document{}, &models.DocumentText{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Init Services
	services.InitVosk()

	// Init Gin
	if config.AppConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Routes
	routes.SetupRoutes(r)

	// Start Server
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
