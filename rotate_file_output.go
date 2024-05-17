package logger

import (
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

// cron
const (
	Hourly = "@hourly"
	Daily  = "@daily"
)

// NewRotateFileOutput new rotate file output
func NewRotateFileOutput(
	filename string,
	maxSize int,
	maxAge int,
	maxBackups int,
	localTime bool,
	compress bool,
	rotationCron string,
) *RotateFileOutput {
	output := new(RotateFileOutput)
	output.Logger = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	}
	if rotationCron != "" {
		cronOptions := []cron.Option{cron.WithSeconds()}
		if localTime {
			cronOptions = append(cronOptions, cron.WithLocation(time.Local))
		} else {
			cronOptions = append(cronOptions, cron.WithLocation(time.UTC))
		}
		output.cron = cron.New(cronOptions...)
		output.cron.AddFunc(rotationCron, func() {
			output.Logger.Rotate()
		})
		output.cron.Start()
	}
	return output
}

// RotateFileOutput rotate file output
type RotateFileOutput struct {
	*lumberjack.Logger
	RotationTime time.Duration
	cron         *cron.Cron
}
