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
	authHandler "github.com/streamverse/auth-service/handlers"
	"github.com/streamverse/auth-service/repository"
	"github.com/streamverse/auth-service/service"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	// Initialize service
	authService := service.NewAuthService(userRepo, tokenRepo, deviceRepo, cfg.JWT.SecretKey)

	// Initialize handlers
	authHandler := authHandler.NewAuthHandler(authService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(nil, 10, 1*time.Minute) // 10 requests per minute

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Auth routes (public) - Issue #11: Routes updated to match requirements
	auth := router.Group("/auth")
	{
		auth.POST("/register", rateLimiter.RateLimit(), authHandler.Register)
		auth.POST("/login", rateLimiter.RateLimit(), authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(cfg.JWT.SecretKey), authHandler.Logout)
		auth.GET("/validate", middleware.AuthMiddleware(cfg.JWT.SecretKey), authHandler.Validate)
		auth.POST("/mfa/setup", middleware.AuthMiddleware(cfg.JWT.SecretKey), authHandler.SetupMFA)
		auth.POST("/mfa/verify", middleware.AuthMiddleware(cfg.JWT.SecretKey), authHandler.VerifyMFA)
		auth.POST("/oauth/google", rateLimiter.RateLimit(), authHandler.OAuthGoogle)
		auth.POST("/oauth/apple", rateLimiter.RateLimit(), authHandler.OAuthApple)
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

	log.Info("Auth service started", logger.String("address", srv.Addr))

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

