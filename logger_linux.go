//go:build linux

package applog

import "log/slog"

func platformHandlers(options loggerOptions) []slog.Handler {
	handlers := []slog.Handler{&JournaldHandler{}}
	if options.verbose {
		handlers = append(handlers, textHandler(outputWriter(options)))
	}
	return handlers
}
