package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/config"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/common-go/middleware"
	transcodingHandler "github.com/streamverse/transcoding-service/handlers"
	"github.com/streamverse/transcoding-service/repository"
	"github.com/streamverse/transcoding-service/service"
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

	// Initialize S3 client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Fatal("Failed to create S3 session", logger.Error(err))
	}
	s3Client := s3.New(sess)

	// Initialize repository
	transcodingRepo := repository.NewTranscodingRepository(db, s3Client)

	// Initialize service
	transcodingService := service.NewTranscodingService(transcodingRepo)

	// Initialize handlers
	transcodingHandler := transcodingHandler.NewTranscodingHandler(transcodingService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Transcoding routes - Issue #15: Routes updated to match requirements
	api := router.Group("/transcode")
	api.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
	{
		api.POST("/jobs", transcodingHandler.SubmitTranscodeJob)        // POST /transcode/jobs
		api.GET("/jobs/:job_id", transcodingHandler.GetTranscodeJobStatus) // GET /transcode/jobs/{job_id}
		api.GET("/jobs", transcodingHandler.ListTranscodeJobs)          // GET /transcode/jobs (with filters)
		api.GET("/profiles", transcodingHandler.ListProfiles)           // GET /transcode/profiles
		api.POST("/profiles", transcodingHandler.CreateProfile)         // POST /transcode/profiles

		// Resumable Upload Routes - Issue #29
		api.POST("/uploads", transcodingHandler.InitiateUpload)
		api.POST("/uploads/:upload_id/parts", transcodingHandler.UploadPart)
		api.POST("/uploads/:upload_id/complete", transcodingHandler.CompleteUpload)
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

	log.Info("Transcoding service started", logger.String("address", srv.Addr))

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
