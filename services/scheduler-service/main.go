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
	schedulerHandler "github.com/streamverse/scheduler-service/handlers"
	"github.com/streamverse/scheduler-service/repository"
	"github.com/streamverse/scheduler-service/service"
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
	schedulerRepo := repository.NewSchedulerRepository(db)

	// Initialize service
	cdnBaseURL := os.Getenv("CDN_BASE_URL")
	if cdnBaseURL == "" {
		cdnBaseURL = "https://cdn.streamverse.io"
	}
	schedulerService := service.NewSchedulerService(schedulerRepo, cdnBaseURL)

	// Initialize handlers
	schedulerHandler := schedulerHandler.NewSchedulerHandler(schedulerService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Scheduler routes - Issue #27: All endpoints public (can add auth if needed)
	api := router.Group("/scheduler")
	{
		api.GET("/channels", schedulerHandler.ListChannels)                              // GET /scheduler/channels
		api.GET("/channels/:channel_id/epg", schedulerHandler.GetChannelEPG)            // GET /scheduler/channels/{channel_id}/epg
		api.GET("/channels/:channel_id/manifest", schedulerHandler.GetChannelManifest)   // GET /scheduler/channels/{channel_id}/manifest
		api.GET("/channels/:channel_id/now", schedulerHandler.GetCurrentScheduleEntry)   // GET /scheduler/channels/{channel_id}/now
		// Admin routes (optional - add auth middleware)
		api.POST("/schedule", schedulerHandler.CreateScheduleEntry)                      // POST /scheduler/schedule
		api.PUT("/schedule/:id", schedulerHandler.UpdateScheduleEntry)                  // PUT /scheduler/schedule/{id}
		api.DELETE("/schedule/:id", schedulerHandler.DeleteScheduleEntry)               // DELETE /scheduler/schedule/{id}
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

	log.Info("Scheduler service started", logger.String("address", srv.Addr))

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
