package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/training-bot-service/internal/bot"
	"github.com/streamverse/training-bot-service/internal/generator"
	"github.com/streamverse/training-bot-service/pkg/ai"
)

func main() {
	// Initialize AI client (Gemini)
	aiClient, err := ai.NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		log.Fatalf("Failed to initialize AI client: %v", err)
	}

	// Initialize training bot
	trainingBot := bot.NewTrainingBot(aiClient)

	// Initialize content generator
	contentGenerator := generator.NewContentGenerator(aiClient)

	// Setup HTTP server
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Training bot chat
		api.POST("/bot/chat", trainingBot.HandleChat)
		api.POST("/bot/session", trainingBot.CreateSession)
		api.GET("/bot/session/:sessionId", trainingBot.GetSession)

		// Content generation
		api.POST("/generate/manual", contentGenerator.GenerateManual)
		api.POST("/generate/video-script", contentGenerator.GenerateVideoScript)
		api.POST("/generate/quiz", contentGenerator.GenerateQuiz)
		api.POST("/generate/onboarding", contentGenerator.GenerateOnboarding)

		// Training content management
		api.GET("/training/manuals", getManuals)
		api.GET("/training/manuals/:category", getManualByCategory)
		api.GET("/training/videos", getVideoTrainings)
		api.GET("/training/modules", getInteractiveModules)

		// User progress tracking
		api.GET("/progress/:userId", getUserProgress)
		api.POST("/progress/:userId/complete", markModuleComplete)

		// Analytics
		api.GET("/analytics/training", getTrainingAnalytics)
	}

	// Setup server
	srv := &http.Server{
		Addr:    ":8096",
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Println("Training Bot Service starting on :8096")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func getManuals(c *gin.Context) {
	// Implementation
	c.JSON(http.StatusOK, gin.H{"manuals": []string{}})
}

func getManualByCategory(c *gin.Context) {
	category := c.Param("category")
	c.JSON(http.StatusOK, gin.H{"category": category})
}

func getVideoTrainings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"videos": []string{}})
}

func getInteractiveModules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"modules": []string{}})
}

func getUserProgress(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "progress": 0})
}

func markModuleComplete(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "status": "completed"})
}

func getTrainingAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"analytics": gin.H{}})
}
