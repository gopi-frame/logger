package slog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSource(t *testing.T) {
	cfg := NewConfig()
	assert.False(t, cfg.AddSource)
	if err := cfg.Apply(AddSource()); err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.True(t, cfg.AddSource)
}

func TestFields(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, make(map[string]any), cfg.Fields)
	if err := cfg.Apply(WithFields(map[string]any{"key": "value"})); err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, map[string]any{"key": "value"}, cfg.Fields)
}

func TestPanicOnFatal(t *testing.T) {
	cfg := NewConfig()
	assert.False(t, cfg.PanicOnFatal)
	if err := cfg.Apply(PanicOnFatal()); err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.True(t, cfg.PanicOnFatal)
}

func TestUnmarshalOptions(t *testing.T) {
	var options = map[string]any{
		"level": "panic",
	}
	cfg, err := UnmarshalOptions(options)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, LevelPanic, cfg.Level.Level)
}
