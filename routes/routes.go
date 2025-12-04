package routes

import (
	"meetingassist/handlers"
	"meetingassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Static files
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	// HTML Templates
	r.LoadHTMLGlob("templates/*")

	// Public Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "landing.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	authHandler := handlers.NewAuthHandler()
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		audioHandler := handlers.NewAudioHandler()
		protected.POST("/audio/upload", audioHandler.UploadAudio)
		protected.GET("/audio", audioHandler.GetAudios)
		protected.GET("/audio/:id", audioHandler.GetAudioDetail)
		protected.PUT("/audio/:id", audioHandler.UpdateTranscript)

		docHandler := handlers.NewDocumentHandler()
		protected.POST("/document/upload", docHandler.UploadDocument)
		protected.GET("/document", docHandler.GetDocuments)
		protected.GET("/document/:id", docHandler.GetDocumentDetail)
		protected.PUT("/document/:id", docHandler.UpdateText)

		exportHandler := handlers.NewExportHandler()
		protected.POST("/export/pdf", exportHandler.ExportPDF)
		protected.POST("/export/docx", exportHandler.ExportDocx)
	}

	// Dashboard Routes (Protected UI)
	// For simplicity, we serve HTML but the frontend will check auth token
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})
	r.GET("/audio/upload", func(c *gin.Context) {
		c.HTML(200, "audio_upload.html", nil)
	})
	r.GET("/document/upload", func(c *gin.Context) {
		c.HTML(200, "document_upload.html", nil)
	})
	r.GET("/history", func(c *gin.Context) {
		c.HTML(200, "history.html", nil)
	})
	r.GET("/profile", func(c *gin.Context) {
		c.HTML(200, "profile.html", nil)
	})
	r.GET("/editor/audio/:id", func(c *gin.Context) {
		c.HTML(200, "audio_editor.html", nil)
	})
	r.GET("/editor/document/:id", func(c *gin.Context) {
		c.HTML(200, "document_editor.html", nil)
	})
}
