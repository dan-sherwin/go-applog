package applog

type Status struct {
	AppName        string
	Level          string
	DebugVerbosity int
	Verbose        bool
	SlogDefault    bool
	DevLogBus      bool
	Generation     uint64
}
