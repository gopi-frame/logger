package logger

import (
	"fmt"
	"strings"

	"github.com/gopi-frame/contract/enum"
	"github.com/gopi-frame/contract/logger"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var allLevels = []Level{
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelPanic,
	LevelFatal,
}

func Levels() []Level {
	return allLevels
}

func (l Level) Enable(level logger.Level) bool {
	return l >= level.(Level)
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
	default:
		return fmt.Sprintf("unknown(%d)", int8(l))
	}
}

func (l Level) Parse(s string) (enum.Enum, error) {
	switch strings.ToLower(s) {
	case "debug", "0":
		return LevelDebug, nil
	case "info", "1":
		return LevelInfo, nil
	case "warn", "warning", "2":
		return LevelWarn, nil
	case "error", "3":
		return LevelError, nil
	case "panic", "5":
		return LevelPanic, nil
	case "fatal", "4":
		return LevelFatal, nil
	default:
		return nil, fmt.Errorf("unknown level %s", s)
	}
}

func (l Level) Equals(other enum.Enum) bool {
	return l == other
}

func (l Level) Values() []enum.Enum {
	values := make([]enum.Enum, len(allLevels))
	for i := range allLevels {
		values[i] = allLevels[i]
	}
	return values
}

func (l Level) Contains(other enum.Enum) bool {
	for _, v := range other.Values() {
		if v == other {
			return true
		}
	}
	return false
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	level, err := new(Level).Parse(string(text))
	if err != nil {
		return err
	}
	*l = level.(Level)
	return nil
}
