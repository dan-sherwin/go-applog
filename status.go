package applog

type Status struct {
	AppName     string
	Level       string
	Verbose     bool
	SlogDefault bool
	DevLogBus   bool
	Generation  uint64
}
