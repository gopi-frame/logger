package logger

import (
	"io"
	"os"
	"sync"

	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/logger"
	events "github.com/gopi-frame/logger/event"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ logger.Logger = (*Logger)(nil)

// NewLogger new logger
func NewLogger(options ...Option) *Logger {
	logger := new(Logger)
	for _, option := range options {
		option(logger)
	}
	return logger
}

// Logger logger
type Logger struct {
	once       sync.Once
	driver     *zap.Logger
	dispatcher event.Dispatcher
	formatter  logger.Formatter
	hooks      []logger.Hook
	outputs    []io.Writer
	zapOptions []zap.Option
}

func (l *Logger) lazyInit() {
	l.once.Do(func() {
		var encoder zapcore.Encoder
		if l.formatter == nil {
			encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey:    "message",
				LevelKey:      "level",
				TimeKey:       "timestamp",
				NameKey:       "name",
				CallerKey:     "function",
				StacktraceKey: "stacktrace",
			})
		} else {
			encoder = newEncoder(l.formatter)
		}
		var writer zapcore.WriteSyncer
		if len(l.outputs) == 0 {
			writer = zapcore.AddSync(os.Stdout)
		} else {
			writers := []zapcore.WriteSyncer{}
			for _, output := range l.outputs {
				writers = append(writers, zapcore.AddSync(output))
			}
			writer = zapcore.NewMultiWriteSyncer(writers...)
		}
		core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
		if len(l.hooks) > 0 {
			for _, hook := range l.hooks {
				core = zapcore.RegisterHooks(core, func(e zapcore.Entry) error {
					if hook.Enable(e.Level.String()) {
						return hook.Handle(&entry{Entry: e})
					}
					return nil
				})
			}
		}
		driver := zap.New(core, l.zapOptions...)
		l.driver = driver
	})
}

// Dispatcher set event dispatcher
func (l *Logger) Dispatcher(d event.Dispatcher) {
	l.dispatcher = d
}

// Formatter set formatter
func (l *Logger) Formatter(formatter logger.Formatter) {
	l.formatter = formatter
}

// Hooks set hook
func (l *Logger) Hooks(hooks ...logger.Hook) {
	l.hooks = hooks
}

// Outputs set outputs
func (l *Logger) Outputs(outputs ...io.Writer) {
	l.outputs = outputs
}

// Debug debug
func (l *Logger) Debug(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Debug(message, values...)
	l.dispatchEvent(zap.DebugLevel.String(), message, fields)
}

// Info info
func (l *Logger) Info(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Info(message, values...)
	l.dispatchEvent(zap.InfoLevel.String(), message, fields)
}

// Warn warn
func (l *Logger) Warn(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Warn(message, values...)
	l.dispatchEvent(zap.WarnLevel.String(), message, fields)
}

// Error error
func (l *Logger) Error(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Error(message, values...)
	l.dispatchEvent(zap.ErrorLevel.String(), message, fields)
}

// Fatal fatal
func (l *Logger) Fatal(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Fatal(message, values...)
	l.dispatchEvent(zap.ErrorLevel.String(), message, fields)
}

// Panic panic
func (l *Logger) Panic(message string, fields map[string]any) {
	l.lazyInit()
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	l.driver.Panic(message, values...)
	l.dispatchEvent(zap.PanicLevel.String(), message, fields)
}

func (l *Logger) dispatchEvent(level string, message string, fields map[string]any) {
	if l.dispatcher != nil {
		l.dispatcher.Dispatch(events.NewMessageLogged(level, message, fields))
	}
}
