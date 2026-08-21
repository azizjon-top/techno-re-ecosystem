package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger wraps zap logger for application use
type Logger struct {
	*zap.Logger
}

var globalLogger *Logger

// Init initializes the global logger
func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.PanicLevel),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	globalLogger = &Logger{logger}
	return nil
}

// Get returns the global logger
func Get() *Logger {
	if globalLogger == nil {
		// Fallback to console logger if not initialized
		logger := zap.NewExample()
		globalLogger = &Logger{logger}
	}
	return globalLogger
}

// GetWithContext returns logger from context or global logger
func GetWithContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value("logger").(*Logger); ok && logger != nil {
		return logger
	}
	return Get()
}

// WithContext returns a new context with the logger attached
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logger", l)
}

// Close closes the logger and flushes buffered logs
func (l *Logger) Close() error {
	return l.Sync()
}

// Convenience functions
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Get().Panic(msg, fields...)
}

// DebugErr logs error at debug level
func DebugErr(err error, msg string, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		Get().Debug(msg, fields...)
	}
}

// ErrorErr logs error at error level
func ErrorErr(err error, msg string, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		Get().Error(msg, fields...)
	}
}

// WarnErr logs error at warn level
func WarnErr(err error, msg string, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		Get().Warn(msg, fields...)
	}
}
