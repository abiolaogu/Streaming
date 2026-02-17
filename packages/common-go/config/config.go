package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const insecureDefaultJWTSecret = "dev-only-change-me"

// Config holds application configuration
type Config struct {
	Server             ServerConfig
	Database           DatabaseConfig
	Redis              RedisConfig
	JWT                JWTConfig
	Logging            LoggingConfig
	ContentServiceAddr string
	PaymentServiceAddr string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI            string
	DatabaseName   string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey              string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level       string
	Development bool
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URI:            getEnv("DATABASE_URI", "mongodb://localhost:27017"),
			DatabaseName:   getEnv("DATABASE_NAME", "streamverse"),
			MaxPoolSize:    getUint64Env("DATABASE_MAX_POOL_SIZE", 100),
			MinPoolSize:    getUint64Env("DATABASE_MIN_POOL_SIZE", 10),
			ConnectTimeout: getDurationEnv("DATABASE_CONNECT_TIMEOUT", 10*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:              getEnv("JWT_SECRET_KEY", insecureDefaultJWTSecret),
			AccessTokenExpiration:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRATION", 15*time.Minute),
			RefreshTokenExpiration: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRATION", 7*24*time.Hour),
		},
		Logging: LoggingConfig{
			Level:       getEnv("LOG_LEVEL", "info"),
			Development: getBoolEnv("LOG_DEVELOPMENT", false),
		},
		ContentServiceAddr: getEnv("CONTENT_SERVICE_ADDR", "localhost:50052"),
		PaymentServiceAddr: getEnv("PAYMENT_SERVICE_ADDR", "localhost:50053"),
	}

	if err := cfg.validate(); err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) validate() error {
	environment := strings.ToLower(strings.TrimSpace(getEnv("ENVIRONMENT", getEnv("APP_ENV", "development"))))
	isProduction := environment == "production" || environment == "prod"
	if !isProduction {
		return nil
	}

	if c.JWT.SecretKey == "" || c.JWT.SecretKey == insecureDefaultJWTSecret {
		return fmt.Errorf("JWT_SECRET_KEY must be configured in production")
	}

	if len(c.JWT.SecretKey) < 32 {
		return fmt.Errorf("JWT_SECRET_KEY must be at least 32 characters in production")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getUint64Env(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return uintValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
