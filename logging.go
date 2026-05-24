package applog

import (
	"context"
	"log/slog"
)

func Debug(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, 1, msg, args...)
}

func Debug2(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, 2, msg, args...)
}

func Debug3(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, 3, msg, args...)
}

func Debug4(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, 4, msg, args...)
}

func Debug5(msg string, args ...any) {
	log(context.Background(), slog.LevelDebug, 5, msg, args...)
}

func Info(msg string, args ...any) {
	log(context.Background(), slog.LevelInfo, 0, msg, args...)
}

func Warn(msg string, args ...any) {
	log(context.Background(), slog.LevelWarn, 0, msg, args...)
}

func Error(msg string, args ...any) {
	log(context.Background(), slog.LevelError, 0, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	log(ctx, level, 0, msg, args...)
}

func log(ctx context.Context, level slog.Level, requiredVerbosity int32, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	runtime := defaultRuntime
	if !runtime.enabled(level) {
		return
	}
	if level == slog.LevelDebug && runtime.debugVerbosity.Load() < requiredVerbosity {
		return
	}
	runtime.currentLogger().Log(ctx, level, msg, args...)
}
