package zap

import (
	"bytes"
	"context"
	"encoding/json"
	loggercontract "github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	DefaultWriter = zapcore.AddSync(buffer)
	var options = map[string]any{
		"level":       "info",
		"development": false,
		"fields":      map[string]any{"key": "value"},
		"caller":      true,
		"callerSkip":  0,
		"encoder":     EncoderJSON,
		"encoderConfig": map[string]any{
			"messageKey":     "message",
			"levelKey":       "level",
			"timeKey":        "time",
			"nameKey":        "name",
			"callerKey":      "caller",
			"functionKey":    "function",
			"stacktraceKey":  "stacktrace",
			"skipLineEnding": false,
			"lineEnding":     "\r\n",
		},
		"fatalHook": "panic",
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
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "info", data["level"])
		assert.Equal(t, "info", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Info")
		assert.Nil(t, data["stacktrace"])
		assert.Empty(t, data["context"])
	})

	t.Run("infof", func(t *testing.T) {
		buffer.Reset()
		l.Infof("infof: %s", "infof")
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "info", data["level"])
		assert.Equal(t, "infof: infof", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Info")
		assert.Nil(t, data["stacktrace"])
		assert.Empty(t, data["context"])
	})

	t.Run("warn", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "warn"})
		l.WithContext(ctx).Warn("warn")
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "warn", data["level"])
		assert.Equal(t, "warn", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Warn")
		assert.Nil(t, data["stacktrace"])
		assert.Equal(t, "warn", data["context"].(map[string]any)["test"])
	})

	t.Run("warnf", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "warn"})
		l.WithContext(ctx).Warnf("warnf: %s", "warnf")
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "warn", data["level"])
		assert.Equal(t, "warnf: warnf", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Warn")
		assert.Nil(t, data["stacktrace"])
		assert.Equal(t, "warn", data["context"].(map[string]any)["test"])
	})

	t.Run("error", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "error"})
		l.WithContext(ctx).Error("error")
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "error", data["level"])
		assert.Equal(t, "error", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Error")
		assert.NotNil(t, data["stacktrace"])
		assert.Equal(t, "error", data["context"].(map[string]any)["test"])
	})

	t.Run("errorf", func(t *testing.T) {
		buffer.Reset()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "error"})
		l.WithContext(ctx).Errorf("errorf: %s", "errorf")
		assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "error", data["level"])
		assert.Equal(t, "errorf: errorf", data["message"])
		assert.Equal(t, "value", data["key"])
		assert.Contains(t, data["caller"], "logger.go")
		assert.Contains(t, data["function"], "Error")
		assert.NotNil(t, data["stacktrace"])
		assert.Equal(t, "error", data["context"].(map[string]any)["test"])
	})

	t.Run("panic", func(t *testing.T) {
		buffer.Reset()
		defer func() {
			if r := recover(); r != nil {
				assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
				data := make(map[string]any)
				err = json.Unmarshal(buffer.Bytes(), &data)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "panic", data["level"])
				assert.Equal(t, "panic", data["message"])
				assert.Equal(t, "value", data["key"])
				assert.Contains(t, data["caller"], "logger.go")
				assert.Contains(t, data["function"], "Panic")
				assert.NotNil(t, data["stacktrace"])
				assert.Equal(t, "panic", data["context"].(map[string]any)["test"])
			}
		}()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "panic"})
		l.WithContext(ctx).Panic("panic")
	})

	t.Run("panicf", func(t *testing.T) {
		buffer.Reset()
		defer func() {
			if r := recover(); r != nil {
				assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
				data := make(map[string]any)
				err = json.Unmarshal(buffer.Bytes(), &data)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "panic", data["level"])
				assert.Equal(t, "panicf: panicf", data["message"])
				assert.Equal(t, "value", data["key"])
				assert.Contains(t, data["caller"], "logger.go")
				assert.Contains(t, data["function"], "Panic")
				assert.NotNil(t, data["stacktrace"])
				assert.Equal(t, "panic", data["context"].(map[string]any)["test"])
			}
		}()
		ctx := logger.WithValue(context.Background(), map[string]any{"test": "panic"})
		l.WithContext(ctx).Panicf("panicf: %s", "panicf")
	})

	t.Run("fatal", func(t *testing.T) {
		buffer.Reset()
		defer func() {
			if r := recover(); r != nil {
				assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
				data := make(map[string]any)
				err = json.Unmarshal(buffer.Bytes(), &data)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "fatal", data["level"])
				assert.Equal(t, "fatal", data["message"])
				assert.Equal(t, "value", data["key"])
				assert.Contains(t, data["caller"], "logger.go")
				assert.Contains(t, data["function"], "Fatal")
				assert.NotNil(t, data["stacktrace"])
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
				assert.True(t, strings.HasSuffix(buffer.String(), "\r\n"))
				data := make(map[string]any)
				err = json.Unmarshal(buffer.Bytes(), &data)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				assert.Equal(t, "fatal", data["level"])
				assert.Equal(t, "fatalf: fatalf", data["message"])
				assert.Equal(t, "value", data["key"])
				assert.Contains(t, data["caller"], "logger.go")
				assert.Contains(t, data["function"], "Fatal")
				assert.NotNil(t, data["stacktrace"])
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
		l.WithLevel(loggercontract.LevelDebug).WithContext(ctx).Debugf("debugf: %s", "debugf")
		assert.NotZero(t, buffer.Len())
		data := make(map[string]any)
		err = json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "debug", data["level"])
		assert.Equal(t, "debugf: debugf", data["message"])
		assert.Equal(t, "debug", data["context"].(map[string]any)["test"])

		buffer.Reset()
		l.Debug("debug")
		assert.Equal(t, "", buffer.String())
	})
}
