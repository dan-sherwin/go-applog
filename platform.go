package applog

import (
	"io"
	"log/slog"
	"os"
)

type loggerOptions struct {
	output             io.Writer
	verbose            bool
	disableSlogDefault bool
	appName            string
	version            string
	commit             string
	buildDate          string
	handlers           []slog.Handler
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
