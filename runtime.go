package applog

import (
	"io"
	"log/slog"
	"os"
	"os/user"
	"strings"
	"sync"
	"sync/atomic"
)

type SetupOptions struct {
	AppName            string
	Version            string
	Commit             string
	BuildDate          string
	Verbose            bool
	Output             io.Writer
	DisableSlogDefault bool
	Handlers           []slog.Handler
}

type runtimeState struct {
	mu                  sync.RWMutex
	appName             string
	version             string
	commit              string
	buildDate           string
	output              io.Writer
	verbose             bool
	disableSlogDefault  bool
	handlers            []slog.Handler
	logger              *slog.Logger
	level               atomic.Int32
	lastConfigurationID atomic.Uint64
}

var defaultRuntime = newRuntimeState()

func newRuntimeState() *runtimeState {
	state := &runtimeState{
		appName: "unknown",
		output:  os.Stdout,
	}
	state.level.Store(rankDebug)
	return state
}

func Setup(options SetupOptions) {
	defaultRuntime.setup(options)
}

func (s *runtimeState) setup(options SetupOptions) {
	appName := strings.TrimSpace(options.AppName)
	if appName == "" {
		appName = "unknown"
	}
	output := options.Output
	if output == nil {
		output = os.Stdout
	}

	s.mu.Lock()
	s.appName = appName
	s.version = strings.TrimSpace(options.Version)
	s.commit = strings.TrimSpace(options.Commit)
	s.buildDate = strings.TrimSpace(options.BuildDate)
	s.output = output
	s.verbose = options.Verbose
	s.disableSlogDefault = options.DisableSlogDefault
	s.handlers = append([]slog.Handler(nil), options.Handlers...)
	s.mu.Unlock()

	s.rebuildLogger()
}

func SetVerbose(verbose bool) {
	defaultRuntime.mu.Lock()
	defaultRuntime.verbose = verbose
	defaultRuntime.mu.Unlock()
	defaultRuntime.rebuildLogger()
}

func Logger() *slog.Logger {
	return defaultRuntime.currentLogger()
}

func With(args ...any) *slog.Logger {
	return Logger().With(args...)
}

func RuntimeStatus() Status {
	return defaultRuntime.status()
}

func (s *runtimeState) currentLogger() *slog.Logger {
	s.mu.RLock()
	logger := s.logger
	s.mu.RUnlock()
	if logger == nil {
		return slog.Default()
	}
	return logger
}

func (s *runtimeState) rebuildLogger() {
	s.mu.RLock()
	options := loggerOptions{
		output:             s.output,
		verbose:            s.verbose,
		disableSlogDefault: s.disableSlogDefault,
		appName:            s.appName,
		version:            s.version,
		commit:             s.commit,
		buildDate:          s.buildDate,
		handlers:           append([]slog.Handler(nil), s.handlers...),
	}
	s.mu.RUnlock()

	handlers := platformHandlers(options)
	for _, handler := range options.handlers {
		if handler != nil {
			handlers = append(handlers, handler)
		}
	}
	if len(handlers) == 0 {
		handlers = append(handlers, slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	handler := &gateHandler{runtime: s, next: combineHandlers(handlers)}
	logger := slog.New(handler).With(s.standardAttrs()...)

	s.mu.Lock()
	s.logger = logger
	s.lastConfigurationID.Add(1)
	s.mu.Unlock()

	if !options.disableSlogDefault {
		slog.SetDefault(logger)
	}
}

func (s *runtimeState) standardAttrs() []any {
	currentUser, err := user.Current()
	username := "unknown"
	if err == nil && currentUser != nil {
		username = currentUser.Username
	}

	s.mu.RLock()
	appName := s.appName
	version := s.version
	commit := s.commit
	buildDate := s.buildDate
	s.mu.RUnlock()

	return []any{
		slog.Int("pid", os.Getpid()),
		slog.String("user", username),
		slog.String("app", appName),
		slog.String("version", version),
		slog.String("commit", commit),
		slog.String("buildDate", buildDate),
	}
}

func (s *runtimeState) enabled(level slog.Level) bool {
	return s.level.Load() >= requiredRankForSlogLevel(level)
}

func (s *runtimeState) status() Status {
	s.mu.RLock()
	appName := s.appName
	verbose := s.verbose
	slogDefault := !s.disableSlogDefault
	s.mu.RUnlock()
	return Status{
		AppName:     appName,
		Level:       levelName(s.level.Load()),
		Verbose:     verbose,
		SlogDefault: slogDefault,
		Generation:  s.lastConfigurationID.Load(),
	}
}
