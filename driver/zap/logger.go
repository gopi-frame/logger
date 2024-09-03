// Package zap is an implementation for [gopi-frame/logger](https://pkg.go.dev/github.com/gopi-frame/logger)
// based on [uber-go/zap](https://github.com/uber-go/zap).
package zap

import (
	"context"
	"fmt"
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[loggercontract.Level]zapcore.Level{
	loggercontract.LevelDebug: zapcore.DebugLevel,
	loggercontract.LevelInfo:  zapcore.InfoLevel,
	loggercontract.LevelWarn:  zapcore.WarnLevel,
	loggercontract.LevelError: zapcore.ErrorLevel,
	loggercontract.LevelFatal: zapcore.FatalLevel,
	loggercontract.LevelPanic: zapcore.PanicLevel,
}

// Logger is a wrapper around [zap.Logger].
// It implements the [loggercontract.Logger] interface.
type Logger struct {
	// Logger is the actual logger that created from the root one by calling [zap.Logger.WithOptions] and [zap.IncreaseLevel]
	*zap.Logger

	// root logger is the logger with min level set to DebugLevel.
	// because [zap.Logger] can only increase the log level, but not decrease it.
	root *zap.Logger

	ctx context.Context
}

// NewLogger creates a new logger
func NewLogger(cfg *Config, opts ...Option) (*Logger, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	if err := cfg.Apply(opts...); err != nil {
		return nil, err
	}
	encoder, err := cfg.ZapEncoder()
	if err != nil {
		return nil, err
	}
	ws, err := cfg.ZapWriters()
	if err != nil {
		return nil, err
	}
	l := new(Logger)
	l.ctx = context.Background()
	core := zapcore.NewCore(encoder, ws, zapcore.DebugLevel)
	l.root = zap.New(core, cfg.ZapOptions()...)
	l.Logger = l.root.WithOptions(zap.IncreaseLevel(cfg.Level))
	return l, nil
}

// WithLevel returns a new logger with the specified level.
func (l *Logger) WithLevel(level loggercontract.Level) loggercontract.Logger {
	lvl := levelMap[level]
	return &Logger{
		ctx:    l.ctx,
		Logger: l.root.WithOptions(zap.IncreaseLevel(lvl)),
		root:   l.root,
	}
}

// WithContext returns a new logger with the specified context.
func (l *Logger) WithContext(ctx context.Context) loggercontract.Logger {
	return &Logger{
		ctx:    ctx,
		Logger: l.Logger,
		root:   l.root,
	}
}

// Debug logs a message at debug level.
func (l *Logger) Debug(message string) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Debug(message, values...)
}

// Debugf logs a formatted message at debug level.
func (l *Logger) Debugf(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Debug(fmt.Sprintf(format, args...), values...)
}

// Info logs a message at info level.
func (l *Logger) Info(message string) {
	var values []zap.Field
	fmt.Println("Value:", logger.GetValue(l.ctx))
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Info(message, values...)
}

// Infof logs a formatted message at info level.
func (l *Logger) Infof(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Info(fmt.Sprintf(format, args...), values...)
}

// Warn logs a message at warn level.
func (l *Logger) Warn(message string) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Warn(message, values...)
}

// Warnf logs a formatted message at warn level.
func (l *Logger) Warnf(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Warn(fmt.Sprintf(format, args...), values...)
}

// Error logs a message at exception level.
func (l *Logger) Error(message string) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Error(message, values...)
}

// Errorf logs a formatted message at exception level.
func (l *Logger) Errorf(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Error(fmt.Sprintf(format, args...), values...)
}

// Panic logs a message at panic level.
func (l *Logger) Panic(message string) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Panic(message, values...)
}

// Panicf logs a formatted message at panic level.
func (l *Logger) Panicf(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Panic(fmt.Sprintf(format, args...), values...)
}

// Fatal logs a message at fatal level.
func (l *Logger) Fatal(message string) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Fatal(message, values...)
}

// Fatalf logs a formatted message at fatal level.
func (l *Logger) Fatalf(format string, args ...any) {
	var values []zap.Field
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, zap.Any("context", value))
	}
	l.Logger.Fatal(fmt.Sprintf(format, args...), values...)
}
