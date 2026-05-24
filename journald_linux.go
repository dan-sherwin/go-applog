//go:build linux

package applog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/coreos/go-systemd/v22/journal"
)

type JournaldHandler struct{}

func (h *JournaldHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *JournaldHandler) Handle(_ context.Context, record slog.Record) error {
	msg := record.Message
	attrs := make(map[string]string)
	record.Attrs(func(attr slog.Attr) bool {
		value := attr.Value.Resolve().String()
		attrs[strings.ToUpper(attr.Key)] = value
		msg += fmt.Sprintf(", %s=%s", attr.Key, value)
		return true
	})

	var priority journal.Priority
	switch record.Level {
	case slog.LevelDebug:
		priority = journal.PriDebug
	case slog.LevelInfo:
		priority = journal.PriInfo
	case slog.LevelWarn:
		priority = journal.PriWarning
	case slog.LevelError:
		priority = journal.PriErr
	default:
		priority = journal.PriNotice
	}
	return journal.Send(msg, priority, attrs)
}

func (h *JournaldHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *JournaldHandler) WithGroup(string) slog.Handler {
	return h
}
