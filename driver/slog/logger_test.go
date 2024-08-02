package slog

import (
	"bytes"
	"context"
	"encoding/json"
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	DefaultWriter = buffer
	var options = map[string]any{
		"level":        "info",
		"encoder":      EncoderJSON,
		"panicOnFatal": true,
	}
	l, err := new(Driver).Open(options)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	t.Run("debug", func(t *testing.T) {
		buffer.Reset()
		l.Debug("debug")
		assert.Equal(t, "", buffer.String())
	})

	t.Run("debugf", func(t *testing.T) {
		buffer.Reset()
		l.Debugf("debugf: %s", "debugf")
		assert.Equal(t, "", buffer.String())
	})

	t.Run("info", func(t *testing.T) {
		buffer.Reset()
		l.Info("info")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "INFO", data[slog.LevelKey])
		assert.Equal(t, "info", data[slog.MessageKey])
	})

	t.Run("infof", func(t *testing.T) {
		buffer.Reset()
		l.Infof("infof: %s", "infof")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "INFO", data[slog.LevelKey])
		assert.Equal(t, "infof: infof", data[slog.MessageKey])
	})

	t.Run("warn", func(t *testing.T) {
		buffer.Reset()
		l.Warn("warn")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "WARN", data[slog.LevelKey])
		assert.Equal(t, "warn", data[slog.MessageKey])
	})

	t.Run("warnf", func(t *testing.T) {
		buffer.Reset()
		l.Warnf("warnf: %s", "warnf")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "WARN", data[slog.LevelKey])
		assert.Equal(t, "warnf: warnf", data[slog.MessageKey])
	})

	t.Run("error", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "error"})
		l.WithContext(ctx).Error("error")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "ERROR", data[slog.LevelKey])
		assert.Equal(t, "error", data[slog.MessageKey])
		assert.Equal(t, "error", data["context"].(map[string]any)["test"])
	})

	t.Run("errorf", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "errorf"})
		l.WithContext(ctx).Errorf("errorf: %s", "errorf")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "ERROR", data[slog.LevelKey])
		assert.Equal(t, "errorf: errorf", data[slog.MessageKey])
		assert.Equal(t, "errorf", data["context"].(map[string]any)["test"])
	})

	t.Run("panic", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "panic"})
		defer func() {
			if r := recover(); r != nil {
				assert.NotZero(t, buffer.Len())
				var data map[string]any
				if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "PANIC", data[slog.LevelKey])
				assert.Equal(t, "panic", data[slog.MessageKey])
				assert.Equal(t, "panic", data["context"].(map[string]any)["test"])
			}
		}()
		l.WithContext(ctx).Panic("panic")
	})

	t.Run("panicf", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "panic"})
		defer func() {
			if r := recover(); r != nil {
				assert.NotZero(t, buffer.Len())
				var data map[string]any
				if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "PANIC", data[slog.LevelKey])
				assert.Equal(t, "panicf: panicf", data[slog.MessageKey])
				assert.Equal(t, "panic", data["context"].(map[string]any)["test"])
			}
		}()
		l.WithContext(ctx).Panicf("panicf: %s", "panicf")
	})

	t.Run("fatal", func(t *testing.T) {
		buffer.Reset()
		defer func() {
			if r := recover(); r != nil {
				assert.NotZero(t, buffer.Len())
				var data map[string]any
				if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "FATAL", data[slog.LevelKey])
				assert.Equal(t, "fatal", data[slog.MessageKey])
				assert.Equal(t, "fatal", data["context"].(map[string]any)["test"])
			}
		}()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "fatal"})
		l.WithContext(ctx).Fatal("fatal")
	})

	t.Run("fatalf", func(t *testing.T) {
		buffer.Reset()
		defer func() {
			if r := recover(); r != nil {
				assert.NotZero(t, buffer.Len())
				var data map[string]any
				if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "FATAL", data[slog.LevelKey])
				assert.Equal(t, "fatalf: fatalf", data[slog.MessageKey])
				assert.Equal(t, "fatal", data["context"].(map[string]any)["test"])
			}
		}()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "fatal"})
		l.WithContext(ctx).Fatalf("fatalf: %s", "fatalf")
	})

	t.Run("decrease level", func(t *testing.T) {
		buffer.Reset()
		l.Debug("debug")
		assert.Equal(t, "", buffer.String())
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "debug"})
		l.WithContext(ctx).WithLevel(loggercontract.LevelDebug).Debugf("debugf: %s", "debugf")
		assert.NotZero(t, buffer.Len())
		var data map[string]any
		if err := json.Unmarshal(buffer.Bytes(), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "DEBUG", data[slog.LevelKey])
		assert.Equal(t, "debugf: debugf", data[slog.MessageKey])
		assert.Equal(t, "debug", data["context"].(map[string]any)["test"])

		buffer.Reset()
		l.Debug("debug")
		assert.Equal(t, "", buffer.String())
	})
}
