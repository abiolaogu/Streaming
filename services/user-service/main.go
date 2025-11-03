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
	userHandler "github.com/streamverse/user-service/handlers"
	"github.com/streamverse/user-service/repository"
	"github.com/streamverse/user-service/service"
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
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := userHandler.NewUserHandler(userService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// User routes
	api := router.Group("/api/v1/users")
	{
		me := api.Group("/me")
		{
			me.GET("/profile", userHandler.GetProfile)
			me.PUT("/profile", userHandler.UpdateProfile)
			me.GET("/preferences", userHandler.GetPreferences)
			me.PUT("/preferences", userHandler.UpdatePreferences)
			me.GET("/profiles", userHandler.GetProfiles)
			me.POST("/profiles", userHandler.CreateProfile)
			me.GET("/watch-history", userHandler.GetWatchHistory)
			me.PUT("/watch-history", userHandler.UpdateWatchHistory)
			me.DELETE("/watch-history", userHandler.ClearWatchHistory)
			me.GET("/continue-watching", userHandler.GetContinueWatching)
			me.GET("/watchlist", userHandler.GetWatchlist)
			me.POST("/watchlist", userHandler.AddToWatchlist)
			me.DELETE("/watchlist/:contentId", userHandler.RemoveFromWatchlist)
			me.DELETE("", userHandler.DeleteUser)
		}
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

	log.Info("User service started", logger.String("address", srv.Addr))

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

