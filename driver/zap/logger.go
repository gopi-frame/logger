package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *Config) (*Logger, error) {
	encoder, err := cfg.ZapEncoder()
	if err != nil {
		return nil, err
	}
	ws, err := cfg.ZapWriters()
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(encoder, ws, cfg.Level)
	logger := zap.New(core, cfg.ZapOptions()...)
	return &Logger{logger}, nil
}

type Logger struct {
	*zap.Logger
}

func (z *Logger) Debug(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	fmt.Println(z.Logger)
	z.Logger.Debug(message, values...)
	z.dispatchEvent(zap.DebugLevel.String(), message, fields)
}

func (z *Logger) Info(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	z.Logger.Info(message, values...)
	z.dispatchEvent(zap.InfoLevel.String(), message, fields)
}

func (z *Logger) Warn(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	z.Logger.Warn(message, values...)
	z.dispatchEvent(zap.WarnLevel.String(), message, fields)
}

func (z *Logger) Error(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	z.Logger.Error(message, values...)
	z.dispatchEvent(zap.ErrorLevel.String(), message, fields)
}

func (z *Logger) Fatal(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	z.Logger.Fatal(message, values...)
	z.dispatchEvent(zap.ErrorLevel.String(), message, fields)
}

func (z *Logger) Panic(message string, fields map[string]any) {
	values := []zap.Field{zap.Namespace("context")}
	for key, value := range fields {
		values = append(values, zap.Any(key, value))
	}
	z.Logger.Panic(message, values...)
	z.dispatchEvent(zap.PanicLevel.String(), message, fields)
}

func (z *Logger) dispatchEvent(level string, message string, fields map[string]any) {

}
