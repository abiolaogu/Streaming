package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a structured logger wrapper
type Logger struct {
	*zap.Logger
}

// New creates a new logger with the specified level
func New(level string, development bool) (*Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewDevelopmentConfig()
	if !development {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// NewDefault creates a logger with default settings
func NewDefault() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{Logger: logger}
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithError adds an error field
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

