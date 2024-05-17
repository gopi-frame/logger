package logger

import (
	"io"

	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/logger"
	"go.uber.org/zap"
)

// Option option
type Option func(*Logger)

// Dispatcher set dispatcher
func Dispatcher(dispatcher event.Dispatcher) Option {
	return func(l *Logger) {
		l.dispatcher = dispatcher
	}
}

// Formatter set formatter
func Formatter(formatter logger.Formatter) Option {
	return func(l *Logger) {
		l.formatter = formatter
	}
}

// Hooks set hooks
func Hooks(hooks ...logger.Hook) Option {
	return func(l *Logger) {
		l.hooks = hooks
	}
}

// Outputs set outputs
func Outputs(outputs ...io.Writer) Option {
	return func(l *Logger) {
		l.outputs = outputs
	}
}

// ZapOptions set zap core options
func ZapOptions(options ...zap.Option) Option {
	return func(l *Logger) {
		l.zapOptions = options
	}
}
