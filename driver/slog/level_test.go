package slog

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestLevel_UnmarshalJSON(t *testing.T) {
	t.Run("standard level string", func(t *testing.T) {
		var jsonStr = `{"level": "debug"}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, slog.LevelDebug, data.Level.Level)
	})

	t.Run("standard level integer", func(t *testing.T) {
		var jsonStr = `{"level": 4}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, slog.LevelWarn, data.Level.Level)
	})

	t.Run("standard level integer string", func(t *testing.T) {
		var jsonStr = `{"level": "4"}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, slog.LevelWarn, data.Level.Level)
	})

	t.Run("custom level string", func(t *testing.T) {
		var jsonStr = `{"level": "panic"}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, LevelPanic, data.Level.Level)
	})

	t.Run("custom level integer", func(t *testing.T) {
		var jsonStr = `{"level": 12}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, LevelPanic, data.Level.Level)
	})

	t.Run("custom level integer string", func(t *testing.T) {
		var jsonStr = `{"level": "12"}`
		var data struct {
			Level Level `json:"level"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, LevelPanic, data.Level.Level)
	})
}
