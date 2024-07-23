package gormlogger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gopi-frame/contract/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	logger                    logger.Logger
	Level                     gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func New(logger logger.Logger) *Logger {
	return &Logger{
		logger:                    logger,
		Level:                     gormlogger.Warn,
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             200 * time.Millisecond,
	}
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &Logger{
		Level:                     level,
		logger:                    l.logger,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		SlowThreshold:             l.SlowThreshold,
	}
}

func (l *Logger) Info(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Info {
		return
	}
	l.logger.Info(fmt.Sprintf(message, args...), nil)
}

func (l *Logger) Warn(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Warn {
		return
	}
	l.logger.Warn(fmt.Sprintf(message, args...), nil)
}

func (l *Logger) Error(ctx context.Context, message string, args ...any) {
	if l.Level < gormlogger.Error {
		return
	}
	l.logger.Error(fmt.Sprintf(message, args...), nil)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.Level <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	if err != nil && l.Level >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		sql, rows := fc()
		l.logger.Error("trace", map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
	} else if l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.Level >= gormlogger.Warn {
		sql, rows := fc()
		l.logger.Warn("slow sql", map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
	} else if l.Level >= gormlogger.Info {
		sql, rows := fc()
		l.logger.Info("trace", map[string]any{"elapsed": elapsed, "rows": rows, "sql": sql})
	}
}
