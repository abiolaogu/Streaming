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
	adminHandler "github.com/streamverse/admin-service/handlers"
	"github.com/streamverse/admin-service/repository"
	"github.com/streamverse/admin-service/service"
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
	adminRepo := repository.NewAdminRepository(db)

	// Initialize service
	adminService := service.NewAdminService(adminRepo)

	// Initialize handlers
	adminHandler := adminHandler.NewAdminHandler(adminService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	
	// All admin routes require authentication
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Admin routes - Issue #21: All endpoints require RBAC
	api := router.Group("/admin")
	{
		// User management
		api.GET("/users", adminHandler.ListUsers)           // GET /admin/users
		api.GET("/users/:id", adminHandler.GetUser)         // GET /admin/users/{id}
		api.PUT("/users/:id", adminHandler.UpdateUser)    // PUT /admin/users/{id}
		api.DELETE("/users/:id", adminHandler.DeleteUser) // DELETE /admin/users/{id}

		// Content management
		api.GET("/content", adminHandler.ListContent)          // GET /admin/content
		api.POST("/content", adminHandler.BulkImportContent)   // POST /admin/content (bulk import)
		api.PUT("/content/:id", adminHandler.UpdateContent)   // PUT /admin/content/{id}
		api.DELETE("/content/:id", adminHandler.DeleteContent) // DELETE /admin/content/{id}

		// Analytics
		api.GET("/analytics", adminHandler.GetDashboardMetrics) // GET /admin/analytics

		// Settings
		api.GET("/settings", adminHandler.GetSystemSettings)   // GET /admin/settings
		api.PUT("/settings", adminHandler.UpdateSystemSettings) // PUT /admin/settings

		// Audit logs
		api.GET("/audit-logs", adminHandler.GetAuditLogs) // GET /admin/audit-logs
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

	log.Info("Admin service started", logger.String("address", srv.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", logger.Error(err))
	}

	log.Info("Server exited")
}

