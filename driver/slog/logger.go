package slog

import (
	"context"
	"fmt"
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
	"log/slog"
	"os"
)

var levelMap = map[loggercontract.Level]slog.Level{
	loggercontract.LevelDebug: slog.LevelDebug,
	loggercontract.LevelInfo:  slog.LevelInfo,
	loggercontract.LevelWarn:  slog.LevelWarn,
	loggercontract.LevelError: slog.LevelError,
	loggercontract.LevelPanic: LevelPanic,
	loggercontract.LevelFatal: LevelFatal,
}

// Logger is a [slog] backed logger.
type Logger struct {
	*slog.Logger

	ctx          context.Context
	level        slog.Leveler
	panicOnFatal bool
}

// NewLogger creates a new logger.
func NewLogger(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	handler, err := cfg.SlogHandler()
	if err != nil {
		return nil, err
	}
	l := slog.New(handler)
	if cfg.Fields != nil {
		var args []any
		for key, value := range cfg.Fields {
			args = append(args, key, value)
		}
		l = l.With(args...)
	}
	return &Logger{
		Logger:       l,
		level:        cfg.Level.Level,
		ctx:          context.WithValue(context.Background(), levelKey, cfg.Level.Level),
		panicOnFatal: cfg.PanicOnFatal,
	}, nil
}

func (l *Logger) WithLevel(level loggercontract.Level) loggercontract.Logger {
	var lvl = l.level
	if value, ok := levelMap[level]; ok {
		lvl = value
	}
	return &Logger{
		Logger:       l.Logger,
		level:        lvl,
		ctx:          context.WithValue(l.ctx, levelKey, lvl),
		panicOnFatal: l.panicOnFatal,
	}
}

func (l *Logger) WithContext(ctx context.Context) loggercontract.Logger {
	return &Logger{
		Logger:       l.Logger,
		level:        l.level,
		ctx:          context.WithValue(ctx, levelKey, l.level),
		panicOnFatal: l.panicOnFatal,
	}
}

// Debug logs a message at [slog.LevelDebug].
func (l *Logger) Debug(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.DebugContext(l.ctx, message, values...)
}

func (l *Logger) Debugf(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.DebugContext(l.ctx, fmt.Sprintf(format, args...), values...)
}

// Info logs a message at [slog.LevelInfo].
func (l *Logger) Info(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.InfoContext(l.ctx, message, values...)
}

func (l *Logger) Infof(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.InfoContext(l.ctx, fmt.Sprintf(format, args...), values...)
}

// Warn logs a message at [slog.LevelWarn].
func (l *Logger) Warn(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.WarnContext(l.ctx, message, values...)
}

func (l *Logger) Warnf(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.WarnContext(l.ctx, fmt.Sprintf(format, args...), values...)
}

// Error logs a message at [slog.LevelError].
func (l *Logger) Error(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.ErrorContext(l.ctx, message, values...)
}

func (l *Logger) Errorf(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.ErrorContext(l.ctx, fmt.Sprintf(format, args...), values...)
}

// Panic logs a message at [LevelPanic].
func (l *Logger) Panic(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.Log(l.ctx, LevelPanic, message, values...)
	panic(message)
}

func (l *Logger) Panicf(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.Log(l.ctx, LevelPanic, fmt.Sprintf(format, args...), values...)
	panic(fmt.Sprintf(format, args...))
}

// Fatal logs a message at [LevelFatal].
func (l *Logger) Fatal(message string) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.Log(l.ctx, LevelFatal, message, values...)
	if l.panicOnFatal {
		panic(message)
	}
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	var values []any
	if value := logger.GetValue(l.ctx); value != nil {
		values = append(values, slog.Any("context", value))
	}
	l.Logger.Log(l.ctx, LevelFatal, fmt.Sprintf(format, args...), values...)
	if l.panicOnFatal {
		panic(fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}
