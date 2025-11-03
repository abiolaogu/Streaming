# Common Go Package

Shared utilities and helpers for Go microservices.

## Features

- ✅ Structured logging with zap
- ✅ Error handling with custom error types
- ✅ HTTP middleware (CORS, auth, rate limiting)
- ✅ MongoDB database helpers
- ✅ JWT utilities
- ✅ Configuration loader
- ✅ Metrics and tracing helpers

## Usage

```go
import (
    "github.com/streamverse/common-go/logger"
    "github.com/streamverse/common-go/middleware"
    "github.com/streamverse/common-go/jwt"
    "github.com/streamverse/common-go/config"
)

// Initialize logger
log, _ := logger.New("info", false)

// Load config
cfg := config.Load()

// Generate JWT token
token, _ := jwt.GenerateAccessToken("user123", "user@example.com", []string{"user"}, cfg.JWT.SecretKey, cfg.JWT.AccessTokenExpiration)

// Use middleware
router.Use(middleware.CORS())
router.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
```

## Documentation

See individual package README files for detailed usage.

