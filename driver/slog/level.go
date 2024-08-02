package slog

import (
	"bytes"
	"log/slog"
	"strconv"
	"strings"
)

// Extra logger levels
const (
	LevelPanic = slog.Level(12)
	LevelFatal = slog.Level(16)
)

// Labels for extra logger levels
const (
	LabelPanic = "panic"
	LabelFatal = "fatal"
)

var levelLabels = map[slog.Leveler]string{
	LevelPanic: LabelPanic,
	LevelFatal: LabelFatal,
}

type Level struct {
	slog.Level
}

func (l *Level) UnmarshalText(text []byte) error {
	s := string(text)
	if i, err := strconv.Atoi(string(text)); err == nil {
		l.Level = slog.Level(i)
		return nil
	}
	if strings.ToLower(s) == LabelPanic {
		l.Level = LevelPanic
		return nil
	}
	if bytes.Equal(bytes.ToLower(text), []byte(LabelFatal)) {
		l.Level = LevelFatal
		return nil
	}
	return l.Level.UnmarshalText(text)
}

func (l *Level) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return l.UnmarshalText(data)
	}
	return l.UnmarshalText([]byte(s))
}

func (l *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return l.UnmarshalText([]byte(s))
}
