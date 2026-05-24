package applog

import (
	"errors"

	app_settings "github.com/dan-sherwin/go-app-settings"
)

type EmptyArgs struct{}

type LevelArgs struct {
	Level string
}

type VerbosityArgs struct {
	Verbosity int
}

type ConfigureArgs struct {
	Level          string
	DebugVerbosity int
}

type rpcReceiver struct {
	runtime *runtimeState
	persist bool
}

func newRPCReceiver(runtime *runtimeState, persist bool) *rpcReceiver {
	return &rpcReceiver{runtime: runtime, persist: persist}
}

func (r *rpcReceiver) Status(_ EmptyArgs, reply *Status) error {
	status, err := r.status()
	if err != nil {
		return err
	}
	*reply = status
	return nil
}

func (r *rpcReceiver) Configure(args ConfigureArgs, reply *Status) error {
	if _, _, err := parseLevel(args.Level); err != nil {
		return err
	}
	if args.DebugVerbosity < 0 || args.DebugVerbosity > 5 {
		return errors.New("debug verbosity must be between 0 and 5")
	}
	if r.persist {
		if err := app_settings.SetSetting(settingLogLevel, args.Level); err != nil {
			return err
		}
		if err := app_settings.SetSetting(settingDebugVerbosity, args.DebugVerbosity); err != nil {
			return err
		}
	} else {
		if err := SetLevel(args.Level); err != nil {
			return err
		}
		if err := SetDebugVerbosity(args.DebugVerbosity); err != nil {
			return err
		}
	}
	return r.Status(EmptyArgs{}, reply)
}

func (r *rpcReceiver) SetLevel(args LevelArgs, reply *Status) error {
	if _, _, err := parseLevel(args.Level); err != nil {
		return err
	}
	if r.persist {
		if err := app_settings.SetSetting(settingLogLevel, args.Level); err != nil {
			return err
		}
	} else if err := SetLevel(args.Level); err != nil {
		return err
	}
	return r.Status(EmptyArgs{}, reply)
}

func (r *rpcReceiver) SetDebugVerbosity(args VerbosityArgs, reply *Status) error {
	if args.Verbosity < 0 || args.Verbosity > 5 {
		return errors.New("debug verbosity must be between 0 and 5")
	}
	if r.persist {
		if err := app_settings.SetSetting(settingDebugVerbosity, args.Verbosity); err != nil {
			return err
		}
	} else if err := SetDebugVerbosity(args.Verbosity); err != nil {
		return err
	}
	return r.Status(EmptyArgs{}, reply)
}

func (r *rpcReceiver) status() (Status, error) {
	if r == nil || r.runtime == nil {
		return Status{}, errors.New("applog runtime is not configured")
	}
	return r.runtime.status(), nil
}

func currentStatus() (Status, error) {
	var status Status
	if err := defaultRuntime.call("Status", EmptyArgs{}, &status); err != nil {
		return Status{}, err
	}
	return status, nil
}

func setRuntimeLevel(level string) (Status, error) {
	var status Status
	if err := defaultRuntime.call("SetLevel", LevelArgs{Level: level}, &status); err != nil {
		return Status{}, err
	}
	return status, nil
}

func setRuntimeDebugVerbosity(verbosity int) (Status, error) {
	var status Status
	if err := defaultRuntime.call("SetDebugVerbosity", VerbosityArgs{Verbosity: verbosity}, &status); err != nil {
		return Status{}, err
	}
	return status, nil
}

func (s *runtimeState) call(method string, args any, reply any) error {
	s.mu.RLock()
	callRPC := s.callRPC
	rpcName := s.rpcName
	s.mu.RUnlock()
	if callRPC == nil {
		return errors.New("applog RPC caller is not configured")
	}
	if rpcName == "" {
		rpcName = defaultRPCName
	}
	return callRPC(rpcName+"."+method, args, reply)
}
