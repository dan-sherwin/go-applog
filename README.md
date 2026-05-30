# GoAppLog

Opinionated application logging for Dan's Go services and CLIs.

`applog` wraps the standard library `slog` package without replacing it. The
package owns only the reusable public logging core:

- platform logger wiring for Linux, macOS, and fallback platforms
- standard fields such as app, version, commit, build date, user, and pid
- helper functions for normal and noisy debug logging
- runtime log-level filtering independent of `slog.LevelVar`

App settings, Kong command trees, RPC receivers, and DevLogBus publishing live
outside this module so public SDK consumers can use `go-applog` without pulling
template/runtime dependencies.

## Usage

Register logging once from the app setup path:

```go
applog.Setup(applog.SetupOptions{
	AppName:   consts.APPNAME,
	Version:   consts.Version,
	Commit:    consts.Commit,
	BuildDate: consts.BuildDate,
})
```

Attach additional `slog.Handler` instances when another module owns an
integration sink:

```go
applog.Setup(applog.SetupOptions{
	AppName:  consts.APPNAME,
	Handlers: extraHandlers,
})
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

## Levels

- `error`
- `warn`
- `info`
- `debug`
- `debug2`
- `debug3`
- `debug4`
- `debug5`

`Debug` requires `log_level=debug` or higher. `Debug2` requires
`log_level=debug2` or higher. `Debug3` requires `log_level=debug3`, and so on
through `Debug5`. The `Debug2` through `Debug5` helpers still emit records at
standard `slog.LevelDebug`; `applog` applies the extra filtering before
emission.
