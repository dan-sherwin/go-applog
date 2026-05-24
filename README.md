# GoAppLog

Opinionated application logging for Dan's Go services.

`applog` wraps the standard library `slog` package without replacing it. The
package owns the repeated service logging setup:

- `log_level` app setting registration
- `debug_verbosity` app setting registration
- platform logger wiring for Linux, macOS, and fallback platforms
- DevLogBus handler setup
- standard fields such as app, version, commit, build date, user, and pid
- Kong command handlers for runtime logging controls
- helper functions for normal and noisy debug logging

The package does not use `slog.LevelVar`. Runtime filtering is controlled by
`applog` before records are emitted.

## Usage

Register logging once from the app setup path:

```go
applog.Setup(applog.SetupOptions{
	AppName:     consts.APPNAME,
	Version:     consts.Version,
	Commit:      consts.Commit,
	BuildDate:   consts.BuildDate,
	RegisterRPC: rpc.RegisterName,
	CallRPC:     rpc.Call,
})
```

After `app_settings.Setup` has loaded saved settings and before Kong parses the
CLI, refresh Kong vars:

```go
applog.ConfigureKongVars(&vars)
```

Embed the commands:

```go
type Commands struct {
	AppLog applog.Commands `cmd:"" name:"logging" help:"Manage runtime application logging"`
}
```

Then log through `applog`:

```go
applog.Info("gRPC reflection disabled")
applog.Debug("Starting signal handler")
applog.Debug2("executing scheduled automation rule", "rule", ruleName)
applog.Debug3("rule evaluation details", "rule", ruleName, "state", state)
```

Plain `slog` remains available. By default `Setup` also configures the process
default `slog` logger so existing `slog.Info` calls continue to use the common
handlers during migration. Set `DisableSlogDefault` when an app wants to keep
raw `slog` completely independent.

## Settings

- `log_level`: `debug`, `info`, `warn`, or `error`
- `debug_verbosity`: `0`, `1`, `2`, or `3`

`Debug` requires `log_level=debug`. `Debug2` also requires
`debug_verbosity >= 2`. `Debug3` requires `debug_verbosity >= 3`.

## Commands

```text
service logging status
service logging level debug
service logging verbosity 2
```
