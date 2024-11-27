package logger

import (
	"context"

	"github.com/gopi-frame/contract/logger"
)

// DeferLogger defer logger
type DeferLogger struct {
	driver string
	config map[string]any

	Logger logger.Logger
}

// NewDeferLogger create a defer logger
func NewDeferLogger(driver string, config map[string]any) *DeferLogger {
	return &DeferLogger{
		driver: driver,
		config: config,
	}
}

func (l *DeferLogger) deferInit() {
	if l.Logger != nil {
		return
	}
	var err error
	if l.Logger, err = Open(l.driver, l.config); err != nil {
		panic(err)
	}
	return
}

func (l *DeferLogger) WithLevel(level logger.Level) logger.Logger {
	l.deferInit()
	return l.Logger.WithLevel(level)
}

func (l *DeferLogger) WithContext(ctx context.Context) logger.Logger {
	l.deferInit()
	return l.Logger.WithContext(ctx)
}

func (l *DeferLogger) Debug(message string) {
	l.deferInit()
	l.Logger.Debug(message)
}

func (l *DeferLogger) Debugf(format string, args ...any) {
	l.deferInit()
	l.Logger.Debugf(format, args...)
}

func (l *DeferLogger) Info(message string) {
	l.deferInit()
	l.Logger.Info(message)
}

func (l *DeferLogger) Infof(format string, args ...any) {
	l.deferInit()
	l.Logger.Infof(format, args...)
}

func (l *DeferLogger) Warn(message string) {
	l.deferInit()
	l.Logger.Warn(message)
}

func (l *DeferLogger) Warnf(format string, args ...any) {
	l.deferInit()
	l.Logger.Warnf(format, args...)
}

func (l *DeferLogger) Error(message string) {
	l.deferInit()
	l.Logger.Error(message)
}

func (l *DeferLogger) Errorf(format string, args ...any) {
	l.deferInit()
	l.Logger.Errorf(format, args...)
}

func (l *DeferLogger) Panic(message string) {
	l.deferInit()
	l.Logger.Panic(message)
}

func (l *DeferLogger) Panicf(format string, args ...any) {
	l.deferInit()
	l.Logger.Panicf(format, args...)
}

func (l *DeferLogger) Fatal(message string) {
	l.deferInit()
	l.Logger.Fatal(message)
}

func (l *DeferLogger) Fatalf(format string, args ...any) {
	l.deferInit()
	l.Logger.Fatalf(format, args...)
}
