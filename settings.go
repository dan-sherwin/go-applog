package applog

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/alecthomas/kong"
	app_settings "github.com/dan-sherwin/go-app-settings"
)

const (
	settingLogLevel       = "log_level"
	settingDebugVerbosity = "debug_verbosity"
	kongLoggingLevelVar   = "logging_level"
)

func ConfigureKongVars(vars *kong.Vars) {
	if vars == nil {
		return
	}
	if *vars == nil {
		*vars = kong.Vars{}
	}
	(*vars)[kongLoggingLevelVar] = Level()
	(*vars)[settingLogLevel] = Level()
	(*vars)[settingDebugVerbosity] = strconv.Itoa(DebugVerbosity())
}

func SetLevel(value string) error {
	level, _, err := parseLevel(value)
	if err != nil {
		return err
	}
	defaultRuntime.level.Store(int32(level))
	return nil
}

func Level() string {
	return levelName(slogLevel())
}

func DebugVerbosity() int {
	return int(defaultRuntime.debugVerbosity.Load())
}

func SetDebugVerbosity(value int) error {
	if value < 0 || value > 5 {
		return fmt.Errorf("debug verbosity must be between 0 and 5")
	}
	defaultRuntime.debugVerbosity.Store(int32(value))
	return nil
}

func registerSettingsHandlers() {
	app_settings.RegisterSetting(&app_settings.Setting{
		Name:        settingLogLevel,
		Description: "Logging level (debug|info|warn|error)",
		GetFunc:     Level,
		SetFunc:     SetLevel,
	})
	app_settings.RegisterSetting(&app_settings.Setting{
		Name:        settingDebugVerbosity,
		Description: "Debug verbosity (0|1|2|3|4|5)",
		GetFunc: func() string {
			return strconv.Itoa(DebugVerbosity())
		},
		SetFunc: func(value string) error {
			i, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			return SetDebugVerbosity(i)
		},
	})
}

func slogLevel() slog.Level {
	return slog.Level(defaultRuntime.level.Load())
}
