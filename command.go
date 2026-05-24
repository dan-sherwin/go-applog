package applog

import (
	"fmt"
	"io"
	"os"

	godevlogbus "github.com/dan-sherwin/go-devlogbus"
)

type CommandDef struct {
	Logging                Commands `cmd:"" name:"logging" help:"Manage runtime application logging" group:"Logging"`
	godevlogbus.CommandDef `group:"Logging"`
}

type Commands struct {
	Status statusCommand `cmd:"" help:"Show runtime application logging status" default:"1"`
	Level  levelCommand  `cmd:"" help:"Set runtime application log level"`
}

type statusCommand struct{}

type levelCommand struct {
	Level string `arg:"" enum:"error,warn,info,debug,debug2,debug3,debug4,debug5" help:"error, warn, info, debug, debug2, debug3, debug4, or debug5" required:""`
}

func (c *statusCommand) Run() error {
	status, err := currentStatus()
	if err != nil {
		status = defaultRuntime.status()
	}
	printStatus(runtimeWriter(), status)
	return err
}

func (c *levelCommand) Run() error {
	status, err := setRuntimeLevel(c.Level)
	if err != nil {
		return err
	}
	printStatus(runtimeWriter(), status)
	return nil
}

func printStatus(writer io.Writer, status Status) {
	if writer == nil {
		writer = io.Discard
	}
	fmt.Fprintf(writer, "App:            %s\n", status.AppName)
	fmt.Fprintf(writer, "Level:          %s\n", status.Level)
	fmt.Fprintf(writer, "Verbose Stdout: %t\n", status.Verbose)
	fmt.Fprintf(writer, "slog Default:   %t\n", status.SlogDefault)
	fmt.Fprintf(writer, "DevLogBus:      %t\n", status.DevLogBus)
	fmt.Fprintf(writer, "Generation:     %d\n", status.Generation)
}

func runtimeWriter() io.Writer {
	defaultRuntime.mu.RLock()
	defer defaultRuntime.mu.RUnlock()
	if defaultRuntime.output == nil {
		return os.Stdout
	}
	return defaultRuntime.output
}
