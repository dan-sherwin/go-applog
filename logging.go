package applog

import (
	"context"
	"log/slog"
)

func Debug(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, rankDebug, msg, args...)
}

func Debug2(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, rankDebug2, msg, args...)
}

func Debug3(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, rankDebug3, msg, args...)
}

func Debug4(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, rankDebug4, msg, args...)
}

func Debug5(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, rankDebug5, msg, args...)
}

func Info(msg string, args ...any) {
	log(context.Background(), slog.LevelInfo, rankInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	log(context.Background(), slog.LevelWarn, rankWarn, msg, args...)
}

func Error(msg string, args ...any) {
	log(context.Background(), slog.LevelError, rankError, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	log(ctx, level, requiredRankForSlogLevel(level), msg, args...)
}

func log(ctx context.Context, level slog.Level, requiredRank int32, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	runtime := defaultRuntime
	if runtime.level.Load() < requiredRank {
		return
	}
	runtime.currentLogger().Log(ctx, level, msg, args...)
}
