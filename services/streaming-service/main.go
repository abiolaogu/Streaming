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
	"github.com/streamverse/streaming-service/internal/clients/content"
	"github.com/streamverse/streaming-service/internal/clients/payment"
	streamingHandler "github.com/streamverse/streaming-service/handlers"
	"github.com/streamverse/streaming-service/repository"
	"github.com/streamverse/streaming-service/service"
)

func main() {
	cfg := config.Load()

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
	streamingRepo := repository.NewStreamingRepository(db)

	// Initialize gRPC clients
	contentClient, err := content.NewClient(cfg.ContentServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to content service", logger.Error(err))
	}
	defer contentClient.Close()

	paymentClient, err := payment.NewClient(cfg.PaymentServiceAddr)
	if err != nil {
		log.Fatal("Failed to connect to payment service", logger.Error(err))
	}
	defer paymentClient.Close()

	// Initialize service
	streamingService := service.NewStreamingService(
		streamingRepo,
		contentClient,
		paymentClient,
		cfg.JWT.SecretKey,
	)

	// Initialize handlers
	streamingHandler := streamingHandler.NewStreamingHandler(streamingService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Streaming routes
	api := router.Group("/streaming")
	{
		// Manifest endpoints with token
		api.GET("/manifest/:content_id/:token.m3u8", streamingHandler.GetHLSManifest)
		api.GET("/manifest/:content_id/:token.mpd", streamingHandler.GetDASHManifest)
		// Token generation
		api.POST("/token", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.GenerateToken)
		// QoE metrics
		api.POST("/qoe", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.SubmitQoE)
		// Session management
		api.POST("/sessions", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.CreateSession)
		api.PUT("/sessions/:sessionId/position", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.UpdatePosition)
		api.POST("/sessions/:sessionId/heartbeat", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.Heartbeat)
		api.DELETE("/sessions/:sessionId", middleware.AuthMiddleware(cfg.JWT.SecretKey), streamingHandler.EndSession)
	}

	// Start server
	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", logger.Error(err))
		}
	}()

	log.Info("Streaming service started", logger.String("address", srv.Addr))

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
