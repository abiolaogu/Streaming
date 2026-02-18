package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/config"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/common-go/middleware"
	"github.com/streamverse/policy-service/handlers"
	"github.com/streamverse/policy-service/service"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Development)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	policyService := service.NewPolicyService()
	policyHandler := handlers.NewPolicyHandler(policyService)

	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	policy := router.Group("/policy/v1")
	policy.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
	{
		policy.POST("/entitlements/evaluate", policyHandler.EvaluateEntitlement)
	}

	srv := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	log.Info("Policy service started", logger.String("address", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start policy service", logger.Error(err))
	}
}
