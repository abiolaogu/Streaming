package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/config"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/common-go/middleware"
	contentHandler "github.com/streamverse/content-service/handlers"
	"github.com/streamverse/content-service/repository"
	"github.com/streamverse/content-service/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Development)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Connect to MongoDB
	db, err := database.Connect(
		cfg.Database.URI,
		cfg.Database.DatabaseName,
		cfg.Database.MaxPoolSize,
		cfg.Database.MinPoolSize,
		cfg.Database.ConnectTimeout,
	)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.Error(err))
	}
	defer db.Disconnect(context.Background())

	// Initialize repository
	contentRepo := repository.NewContentRepository(db)

	// Initialize service
	contentService := service.NewContentService(contentRepo)

	// Initialize handlers
	contentHandler := contentHandler.NewContentHandler(contentService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Content routes - Issue #13: Routes updated to match requirements
	api := router.Group("/content")
	{
		api.GET("/:id", contentHandler.GetContentByID)                     // GET /content/{id}
		api.GET("/search", contentHandler.SearchContent)                  // GET /content/search
		api.GET("/categories", contentHandler.GetCategories)               // GET /content/categories
		api.GET("/trending", contentHandler.GetTrending)                   // GET /content/trending
		api.POST("/:id/ratings", contentHandler.RateContent)              // POST /content/{id}/ratings
		api.GET("/:id/ratings", contentHandler.GetRatings)                 // GET /content/{id}/ratings
		api.GET("/:id/similar", contentHandler.GetSimilar)                 // GET /content/{id}/similar
		api.GET("/:id/entitlements", contentHandler.GetEntitlements)       // GET /content/{id}/entitlements
		// Admin endpoints (optional for Issue #13)
		api.POST("", middleware.AuthMiddleware(cfg.JWT.SecretKey), contentHandler.CreateContent)
		api.PUT("/:id", middleware.AuthMiddleware(cfg.JWT.SecretKey), contentHandler.UpdateContent)
		api.DELETE("/:id", middleware.AuthMiddleware(cfg.JWT.SecretKey), contentHandler.DeleteContent)
	}

	// Start server
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", logger.Error(err))
		}
	}()

	log.Info("Content service started", logger.String("address", srv.Addr))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", logger.Error(err))
	}

	log.Info("Server exited")
}

