package slog

import (
	"context"
	"log/slog"
)

var levelKey = struct {
	name string
}{
	name: "level",
}

type handler struct {
	handler slog.Handler
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	dynamicLevel := ctx.Value(levelKey)
	if dynamicLevel == nil {
		return h.handler.Enabled(ctx, level)
	}
	return level >= dynamicLevel.(slog.Level)
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		handler: h.handler.WithGroup(name),
	}
}
