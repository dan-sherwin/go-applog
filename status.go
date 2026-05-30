package applog

type Status struct {
	AppName     string
	Level       string
	Verbose     bool
	SlogDefault bool
	Generation  uint64
}
