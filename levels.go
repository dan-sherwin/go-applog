package applog

import (
	"fmt"
	"log/slog"
	"strings"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func parseLevel(value string) (slog.Level, string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case LevelDebug:
		return slog.LevelDebug, LevelDebug, nil
	case LevelInfo, "":
		return slog.LevelInfo, LevelInfo, nil
	case LevelWarn, "warning":
		return slog.LevelWarn, LevelWarn, nil
	case LevelError:
		return slog.LevelError, LevelError, nil
	default:
		return slog.LevelInfo, "", fmt.Errorf("unsupported log level %q", value)
	}
}

func levelName(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return LevelDebug
	case slog.LevelInfo:
		return LevelInfo
	case slog.LevelWarn:
		return LevelWarn
	case slog.LevelError:
		return LevelError
	default:
		return level.String()
	}
}
