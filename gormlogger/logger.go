// Package gormlogger is an implementation of [gormlogger.Interface]
package gormlogger

import (
	"context"
	"errors"
	"time"

	"github.com/gopi-frame/logger"

	loggercontract "github.com/gopi-frame/contract/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var levelMap = map[gormlogger.LogLevel]loggercontract.Level{
	gormlogger.Error: logger.LevelError,
	gormlogger.Warn:  logger.LevelWarn,
	gormlogger.Info:  logger.LevelInfo,
}

type Logger struct {
	logger                    loggercontract.Logger
	Level                     gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func New(logger loggercontract.Logger) *Logger {
	return &Logger{
		logger:                    logger,
		Level:                     gormlogger.Warn,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             200 * time.Millisecond,
	}
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l2 := &Logger{
		Level:                     level,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		SlowThreshold:             l.SlowThreshold,
		logger:                    l.logger,
	}
	if value, ok := levelMap[level]; ok {
		l2.logger = l.logger.WithLevel(value)
	}
	return l2
}

func (l *Logger) Info(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Info {
		return
	}
	l.logger.WithContext(ctx).Infof(message, args...)
}

func (l *Logger) Warn(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Warn {
		return
	}
	l.logger.WithContext(ctx).Warnf(message, args...)
}

func (l *Logger) Error(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Error {
		return
	}
	l.logger.WithContext(ctx).Errorf(message, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.Level <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	if err != nil && l.Level >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		sql, rows := fc()
		ctx := logger.WithValue(ctx, map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
		l.logger.WithContext(ctx).Error("trace")
	} else if l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.Level >= gormlogger.Warn {
		sql, rows := fc()
		ctx := logger.WithValue(ctx, map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
		l.logger.WithContext(ctx).Warn("slow sql")
	} else if l.Level >= gormlogger.Info {
		sql, rows := fc()
		ctx := logger.WithValue(ctx, map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
		l.logger.WithContext(ctx).Info("trace")
	}
}
