package applog

import (
	"fmt"
	"log/slog"
	"strings"
)

const (
	LevelError  = "error"
	LevelWarn   = "warn"
	LevelInfo   = "info"
	LevelDebug  = "debug"
	LevelDebug2 = "debug2"
	LevelDebug3 = "debug3"
	LevelDebug4 = "debug4"
	LevelDebug5 = "debug5"
)

const (
	rankError int32 = iota
	rankWarn
	rankInfo
	rankDebug
	rankDebug2
	rankDebug3
	rankDebug4
	rankDebug5
)

type parsedLevel struct {
	rank int32
}

func parseLevel(value string) (parsedLevel, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case LevelError:
		return parsedLevel{rank: rankError}, nil
	case LevelWarn, "warning":
		return parsedLevel{rank: rankWarn}, nil
	case LevelInfo, "":
		return parsedLevel{rank: rankInfo}, nil
	case LevelDebug:
		return parsedLevel{rank: rankDebug}, nil
	case LevelDebug2:
		return parsedLevel{rank: rankDebug2}, nil
	case LevelDebug3:
		return parsedLevel{rank: rankDebug3}, nil
	case LevelDebug4:
		return parsedLevel{rank: rankDebug4}, nil
	case LevelDebug5:
		return parsedLevel{rank: rankDebug5}, nil
	default:
		return parsedLevel{}, fmt.Errorf("unsupported log level %q", value)
	}
}

func levelName(rank int32) string {
	switch rank {
	case rankError:
		return LevelError
	case rankWarn:
		return LevelWarn
	case rankInfo:
		return LevelInfo
	case rankDebug:
		return LevelDebug
	case rankDebug2:
		return LevelDebug2
	case rankDebug3:
		return LevelDebug3
	case rankDebug4:
		return LevelDebug4
	case rankDebug5:
		return LevelDebug5
	default:
		return LevelInfo
	}
}

func requiredRankForSlogLevel(level slog.Level) int32 {
	switch {
	case level >= slog.LevelError:
		return rankError
	case level >= slog.LevelWarn:
		return rankWarn
	case level >= slog.LevelInfo:
		return rankInfo
	default:
		return rankDebug
	}
}
