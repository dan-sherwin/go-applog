package applog

import (
	"io"
	"log/slog"
	"os"
	"time"
)

type loggerOptions struct {
	output             io.Writer
	verbose            bool
	disableSlogDefault bool
	disableDevLogBus   bool
	appName            string
	version            string
	commit             string
	buildDate          string
	queueSize          int
	publishTimeout     time.Duration
}

func outputWriter(options loggerOptions) io.Writer {
	if options.output != nil {
		return options.output
	}
	return os.Stdout
}

func textHandler(writer io.Writer) slog.Handler {
	if writer == nil {
		writer = os.Stdout
	}
	return slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})
}
