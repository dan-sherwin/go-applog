package applog

const SettingLogLevel = "log_level"

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
