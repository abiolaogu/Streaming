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
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/common-go/middleware"
	searchHandler "github.com/streamverse/search-service/handlers"
	"github.com/streamverse/search-service/repository"
	"github.com/streamverse/search-service/service"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Development)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Initialize Elasticsearch client
	esClient, err := repository.NewElasticsearchClient()
	if err != nil {
		log.Fatal("Failed to connect to Elasticsearch", logger.Error(err))
	}

	// Initialize repository
	searchRepo := repository.NewSearchRepository(esClient)

	// Initialize service
	searchService := service.NewSearchService(searchRepo, log)

	// Initialize handlers
	searchHandler := searchHandler.NewSearchHandler(searchService, log)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Search routes - Issue #17: Routes updated to match requirements
	api := router.Group("/search")
	{
		api.GET("/", searchHandler.Search)           // GET /search (with query, filters, pagination, sort)
		api.GET("/suggest", searchHandler.Suggest)   // GET /search/suggest (autocomplete)
		api.GET("/filters", searchHandler.GetFilters) // GET /search/filters (available filters)
		// Admin endpoint
		api.POST("/index", middleware.AuthMiddleware(cfg.JWT.SecretKey), searchHandler.IndexContent)
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

	log.Info("Search service started", logger.String("address", srv.Addr))

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

