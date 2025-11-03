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
	paymentHandler "github.com/streamverse/payment-service/handlers"
	"github.com/streamverse/payment-service/repository"
	"github.com/streamverse/payment-service/service"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Development)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

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

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := paymentHandler.NewPaymentHandler(paymentService, log)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Payment routes - Issue #16: Routes updated to match requirements
	api := router.Group("/payments")
	{
		// Authenticated routes
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		{
			auth.POST("/subscribe", paymentHandler.CreateSubscription)                      // POST /payments/subscribe
			auth.POST("/subscribe/:subscription_id/cancel", paymentHandler.CancelSubscription) // POST /payments/subscribe/{subscription_id}/cancel
			auth.POST("/purchase", paymentHandler.PurchaseContent)                         // POST /payments/purchase
			auth.GET("/entitlements/:user_id", paymentHandler.GetUserEntitlements)           // GET /payments/entitlements/{user_id}
			auth.GET("/plans", paymentHandler.ListPlans)                                    // GET /payments/plans
		}
		// Webhook endpoint (no auth required - Stripe signs the request)
		api.POST("/webhook", paymentHandler.HandleStripeWebhook) // POST /payments/webhook
	}

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

	log.Info("Payment service started", logger.String("address", srv.Addr))

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

