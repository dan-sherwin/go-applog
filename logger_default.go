//go:build !darwin && !linux

package applog

import "log/slog"

func platformHandlers(options loggerOptions) []slog.Handler {
	return []slog.Handler{textHandler(outputWriter(options))}
}
