package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type entry struct {
	zapcore.Entry
	fields map[string]any
}

func (e *entry) Timestamp() time.Time {
	return e.Entry.Time
}

func (e *entry) File() string {
	return e.Entry.Caller.File
}

func (e *entry) Line() int {
	return e.Entry.Caller.Line
}

func (e *entry) Level() string {
	return e.Entry.Level.CapitalString()
}

func (e *entry) Message() string {
	return e.Entry.Message
}

func (e *entry) Fields() map[string]any {
	return e.fields
}
