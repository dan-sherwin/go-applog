# GoAppLog

Opinionated application logging for Dan's Go services.

`applog` wraps the standard library `slog` package without replacing it. The
package owns the repeated service logging setup:

- `log_level` app setting registration
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
	applog.CommandDef
}
```

That exposes both command groups:

```text
service logging status
service logging level debug3
service devlogbus status
service devlogbus enable
service devlogbus disable
service devlogbus setEndpoint /tmp/devlogbus/devlogbus.sock
```

Then log through `applog`:

```go
applog.Info("gRPC reflection disabled")
applog.Debug("Starting signal handler")
applog.Debug2("executing scheduled automation rule", "rule", ruleName)
applog.Debug3("rule evaluation details", "rule", ruleName, "state", state)
applog.Debug4("deep loop detail", "iteration", i)
applog.Debug5("painfully detailed trace", "state", state)
applog.Error("database ping failed", applog.String("error", err.Error()))
```

Plain `slog` remains available. By default `Setup` also configures the process
default `slog` logger so existing `slog.Info` calls continue to use the common
handlers during migration. Set `DisableSlogDefault` when an app wants to keep
raw `slog` completely independent.

## Settings

- `log_level`: `error`, `warn`, `info`, `debug`, `debug2`, `debug3`, `debug4`, or `debug5`

`Debug` requires `log_level=debug` or higher. `Debug2` requires
`log_level=debug2` or higher. `Debug3` requires `log_level=debug3`, and so on
through `Debug5`. The `Debug2` through `Debug5` helpers still emit records at
standard `slog.LevelDebug`; `applog` applies the extra filtering before emission.

## Application Logging Commands

```text
service logging status
service logging level debug3
```
