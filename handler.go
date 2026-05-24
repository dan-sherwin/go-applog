package applog

import (
	"context"
	"log/slog"
)

type TeeHandler struct {
	handlers []slog.Handler
}

func combineHandlers(handlers []slog.Handler) slog.Handler {
	if len(handlers) == 1 {
		return handlers[0]
	}
	return &TeeHandler{handlers: handlers}
}

func (h *TeeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *TeeHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *TeeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &TeeHandler{handlers: newHandlers}
}

func (h *TeeHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &TeeHandler{handlers: newHandlers}
}

type gateHandler struct {
	runtime *runtimeState
	next    slog.Handler
}

func (h *gateHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.runtime != nil && !h.runtime.enabled(level) {
		return false
	}
	return h.next.Enabled(ctx, level)
}

func (h *gateHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.runtime != nil && !h.runtime.enabled(record.Level) {
		return nil
	}
	return h.next.Handle(ctx, record)
}

func (h *gateHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &gateHandler{runtime: h.runtime, next: h.next.WithAttrs(attrs)}
}

func (h *gateHandler) WithGroup(name string) slog.Handler {
	return &gateHandler{runtime: h.runtime, next: h.next.WithGroup(name)}
}
