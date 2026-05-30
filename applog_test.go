package applog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestApplogLevelFiltering(t *testing.T) {
	resetForTest()
	var buf bytes.Buffer
	Setup(SetupOptions{
		AppName:            "test_app",
		Output:             &buf,
		DisableSlogDefault: true,
	})

	if err := SetLevel(LevelInfo); err != nil {
		t.Fatal(err)
	}
	Debug("hidden debug")
	Info("visible info")
	if strings.Contains(buf.String(), "hidden debug") {
		t.Fatalf("debug log was emitted at info level: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "visible info") {
		t.Fatalf("info log was not emitted: %s", buf.String())
	}

	buf.Reset()
	if err := SetLevel(LevelDebug2); err != nil {
		t.Fatal(err)
	}
	Debug("visible debug")
	Debug2("visible debug2", slog.String("rule", "test"))
	Debug3("hidden debug3")
	Debug4("hidden debug4")
	Debug5("hidden debug5")
	out := buf.String()
	if !strings.Contains(out, "visible debug") {
		t.Fatalf("debug log was not emitted: %s", out)
	}
	if !strings.Contains(out, "visible debug2") {
		t.Fatalf("debug2 log was not emitted: %s", out)
	}
	if strings.Contains(out, "hidden debug3") {
		t.Fatalf("debug3 log was emitted below debug3 level: %s", out)
	}
	if strings.Contains(out, "hidden debug4") {
		t.Fatalf("debug4 log was emitted below debug4 level: %s", out)
	}
	if strings.Contains(out, "hidden debug5") {
		t.Fatalf("debug5 log was emitted below debug5 level: %s", out)
	}

	buf.Reset()
	if err := SetLevel(LevelDebug5); err != nil {
		t.Fatal(err)
	}
	Debug4("visible debug4")
	Debug5("visible debug5")
	out = buf.String()
	if !strings.Contains(out, "visible debug4") {
		t.Fatalf("debug4 log was not emitted at debug5 level: %s", out)
	}
	if !strings.Contains(out, "visible debug5") {
		t.Fatalf("debug5 log was not emitted at debug5 level: %s", out)
	}
}

func TestSetupAddsHandlers(t *testing.T) {
	resetForTest()
	var buf bytes.Buffer
	Setup(SetupOptions{
		AppName:            "test_app",
		DisableSlogDefault: true,
		Handlers: []slog.Handler{
			slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
		},
	})

	if err := SetLevel(LevelInfo); err != nil {
		t.Fatal(err)
	}
	Info("extra handler message")
	if !strings.Contains(buf.String(), "extra handler message") {
		t.Fatalf("extra handler did not receive log output: %s", buf.String())
	}
}

func resetForTest() {
	defaultRuntime = newRuntimeState()
}
