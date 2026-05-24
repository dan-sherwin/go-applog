package applog

import (
	"github.com/alecthomas/kong"
	app_settings "github.com/dan-sherwin/go-app-settings"
)

const (
	settingLogLevel     = "log_level"
	kongLoggingLevelVar = "logging_level"
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
}

func SetLevel(value string) error {
	level, err := parseLevel(value)
	if err != nil {
		return err
	}
	defaultRuntime.level.Store(level.rank)
	return nil
}

func Level() string {
	return levelName(defaultRuntime.level.Load())
}

func registerSettingsHandlers() {
	app_settings.RegisterSetting(&app_settings.Setting{
		Name:        settingLogLevel,
		Description: "Logging level (error|warn|info|debug|debug2|debug3|debug4|debug5)",
		GetFunc:     Level,
		SetFunc:     SetLevel,
	})
}
