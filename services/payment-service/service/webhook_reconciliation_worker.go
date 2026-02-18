package service

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streamverse/common-go/logger"
)

// WebhookReconciliationConfig controls periodic replay of failed webhook events.
type WebhookReconciliationConfig struct {
	Enabled        bool
	Interval       time.Duration
	BatchSize      int
	Region         string
	ResidencyMode  string
	AllowedRegions map[string]struct{}
}

// WebhookReconciliationConfigFromEnv builds worker config from environment variables.
func WebhookReconciliationConfigFromEnv() WebhookReconciliationConfig {
	cfg := WebhookReconciliationConfig{
		Enabled:        strings.EqualFold(strings.TrimSpace(os.Getenv("STRIPE_RECON_WORKER_ENABLED")), "true"),
		Interval:       60 * time.Second,
		BatchSize:      50,
		Region:         strings.ToLower(strings.TrimSpace(os.Getenv("PLATFORM_REGION"))),
		ResidencyMode:  strings.ToLower(strings.TrimSpace(os.Getenv("DATA_RESIDENCY_MODE"))),
		AllowedRegions: parseAllowedRegions(os.Getenv("DATA_RESIDENCY_ALLOWED_REGIONS")),
	}

	if cfg.Region == "" {
		cfg.Region = "global"
	}
	if cfg.ResidencyMode == "" {
		cfg.ResidencyMode = "global"
	}
	if v := strings.TrimSpace(os.Getenv("STRIPE_RECON_WORKER_INTERVAL")); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil && parsed > 0 {
			cfg.Interval = parsed
		}
	}
	if v := strings.TrimSpace(os.Getenv("STRIPE_RECON_WORKER_BATCH_SIZE")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cfg.BatchSize = parsed
		}
	}

	return cfg
}

func parseAllowedRegions(raw string) map[string]struct{} {
	regions := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		region := strings.ToLower(strings.TrimSpace(part))
		if region == "" {
			continue
		}
		regions[region] = struct{}{}
	}
	return regions
}

func (c WebhookReconciliationConfig) residencyAllowsRegion() bool {
	if c.ResidencyMode != "strict" {
		return true
	}
	if len(c.AllowedRegions) == 0 {
		return true
	}
	_, ok := c.AllowedRegions[c.Region]
	return ok
}

// WebhookReconciliationWorker periodically replays failed webhook events.
type WebhookReconciliationWorker struct {
	paymentService *PaymentService
	logger         *logger.Logger
	config         WebhookReconciliationConfig
}

// NewWebhookReconciliationWorker creates a reconciliation worker.
func NewWebhookReconciliationWorker(
	paymentService *PaymentService,
	log *logger.Logger,
	config WebhookReconciliationConfig,
) *WebhookReconciliationWorker {
	return &WebhookReconciliationWorker{
		paymentService: paymentService,
		logger:         log,
		config:         config,
	}
}

// Start runs the reconciliation loop until context cancellation.
func (w *WebhookReconciliationWorker) Start(ctx context.Context) {
	if !w.config.Enabled {
		w.logger.Info("Stripe reconciliation worker disabled")
		return
	}
	if !w.config.residencyAllowsRegion() {
		w.logger.Info(
			"Stripe reconciliation worker disabled by residency controls",
			logger.String("region", w.config.Region),
			logger.String("mode", w.config.ResidencyMode),
		)
		return
	}

	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Stripe reconciliation worker stopped")
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *WebhookReconciliationWorker) runOnce(ctx context.Context) {
	replayed, err := w.paymentService.ReplayFailedWebhookEvents(ctx, w.config.BatchSize)
	if err != nil {
		w.logger.Error("Stripe reconciliation replay failed", logger.Error(err))
		return
	}
	if replayed > 0 {
		w.logger.Info("Stripe reconciliation replayed failed webhook events")
	}
}
