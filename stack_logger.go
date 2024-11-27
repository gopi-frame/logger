package logger

import (
	"context"

	"github.com/gopi-frame/contract/logger"
)

type StackLogger struct {
	channels []logger.Logger
}

func NewStackLogger(channels ...logger.Logger) *StackLogger {
	return &StackLogger{channels: channels}
}

func (s *StackLogger) WithLevel(level logger.Level) logger.Logger {
	l := &StackLogger{}
	for _, channel := range s.channels {
		l.channels = append(l.channels, channel.WithLevel(level))
	}
	return l
}

func (s *StackLogger) WithContext(ctx context.Context) logger.Logger {
	l := &StackLogger{}
	for _, channel := range s.channels {
		l.channels = append(l.channels, channel.WithContext(ctx))
	}
	return l
}

func (s *StackLogger) Debug(message string) {
	for _, channel := range s.channels {
		channel.Debug(message)
	}
}

func (s *StackLogger) Debugf(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Debugf(format, args...)
	}
}

func (s *StackLogger) Info(message string) {
	for _, channel := range s.channels {
		channel.Info(message)
	}
}

func (s *StackLogger) Infof(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Infof(format, args...)
	}
}

func (s *StackLogger) Warn(message string) {
	for _, channel := range s.channels {
		channel.Warn(message)
	}
}

func (s *StackLogger) Warnf(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Warnf(format, args...)
	}
}

func (s *StackLogger) Error(message string) {
	for _, channel := range s.channels {
		channel.Error(message)
	}
}

func (s *StackLogger) Errorf(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Errorf(format, args...)
	}
}

func (s *StackLogger) Panic(message string) {
	for _, channel := range s.channels {
		channel.Panic(message)
	}
}

func (s *StackLogger) Panicf(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Panicf(format, args...)
	}
}

func (s *StackLogger) Fatal(message string) {
	for _, channel := range s.channels {
		channel.Fatal(message)
	}
}

func (s *StackLogger) Fatalf(format string, args ...any) {
	for _, channel := range s.channels {
		channel.Fatalf(format, args...)
	}
}
