package applog

import (
	"io"
	"log/slog"
	"os"
	"os/user"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	godevlogbus "github.com/dan-sherwin/go-devlogbus"
)

const defaultRPCName = "AppLog"

type RegisterRPCFunc = godevlogbus.RegisterRPCFunc
type CallRPCFunc = godevlogbus.CallRPCFunc

type SetupOptions struct {
	AppName               string
	Version               string
	Commit                string
	BuildDate             string
	Verbose               bool
	Output                io.Writer
	DisableSlogDefault    bool
	DisableDevLogBus      bool
	RPCName               string
	DevLogBusRPCName      string
	RegisterRPC           RegisterRPCFunc
	CallRPC               CallRPCFunc
	QueueSize             int
	PublishTimeout        time.Duration
	DisableRPCPersistence bool
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
	disableDevLogBus    bool
	rpcName             string
	callRPC             CallRPCFunc
	queueSize           int
	publishTimeout      time.Duration
	settingsRegistered  bool
	rpcRegistered       bool
	devLogBusSetup      bool
	logger              *slog.Logger
	level               atomic.Int32
	lastConfigurationID atomic.Uint64
}

var defaultRuntime = newRuntimeState()

func newRuntimeState() *runtimeState {
	state := &runtimeState{
		appName: "unknown",
		output:  os.Stdout,
		rpcName: defaultRPCName,
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
	rpcName := strings.TrimSpace(options.RPCName)
	if rpcName == "" {
		rpcName = defaultRPCName
	}

	s.mu.Lock()
	s.appName = appName
	s.version = strings.TrimSpace(options.Version)
	s.commit = strings.TrimSpace(options.Commit)
	s.buildDate = strings.TrimSpace(options.BuildDate)
	s.output = output
	s.verbose = options.Verbose
	s.disableSlogDefault = options.DisableSlogDefault
	s.disableDevLogBus = options.DisableDevLogBus
	s.rpcName = rpcName
	s.callRPC = options.CallRPC
	s.queueSize = options.QueueSize
	s.publishTimeout = options.PublishTimeout
	registerSettings := !s.settingsRegistered
	if registerSettings {
		s.settingsRegistered = true
	}
	registerRPC := options.RegisterRPC != nil && !s.rpcRegistered
	if registerRPC {
		s.rpcRegistered = true
	}
	setupDevLogBus := !options.DisableDevLogBus && !s.devLogBusSetup
	if setupDevLogBus {
		s.devLogBusSetup = true
	}
	s.mu.Unlock()

	if registerSettings {
		registerSettingsHandlers()
	}
	if setupDevLogBus {
		godevlogbus.Setup(godevlogbus.SetupOptions{
			Source:                appName,
			RPCName:               options.DevLogBusRPCName,
			RegisterRPC:           options.RegisterRPC,
			CallRPC:               options.CallRPC,
			Output:                output,
			QueueSize:             options.QueueSize,
			PublishTimeout:        options.PublishTimeout,
			DisableRPCPersistence: options.DisableRPCPersistence,
		})
	}
	if registerRPC {
		options.RegisterRPC(rpcName, newRPCReceiver(s, !options.DisableRPCPersistence))
	}
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
		disableDevLogBus:   s.disableDevLogBus,
		appName:            s.appName,
		version:            s.version,
		commit:             s.commit,
		buildDate:          s.buildDate,
		queueSize:          s.queueSize,
		publishTimeout:     s.publishTimeout,
	}
	s.mu.RUnlock()

	handlers := platformHandlers(options)
	if !options.disableDevLogBus {
		handlers = godevlogbus.WithHandler(handlers, slog.LevelDebug)
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
	devLogBus := !s.disableDevLogBus
	s.mu.RUnlock()
	return Status{
		AppName:     appName,
		Level:       levelName(s.level.Load()),
		Verbose:     verbose,
		SlogDefault: slogDefault,
		DevLogBus:   devLogBus,
		Generation:  s.lastConfigurationID.Load(),
	}
}
