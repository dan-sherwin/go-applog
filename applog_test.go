package applog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestApplogLevelAndVerbosityFiltering(t *testing.T) {
	resetForTest()
	var buf bytes.Buffer
	Setup(SetupOptions{
		AppName:            "test_app",
		Output:             &buf,
		DisableDevLogBus:   true,
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
	if err := SetLevel(LevelDebug); err != nil {
		t.Fatal(err)
	}
	if err := SetDebugVerbosity(2); err != nil {
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
		t.Fatalf("debug3 log was emitted below verbosity 3: %s", out)
	}
	if strings.Contains(out, "hidden debug4") {
		t.Fatalf("debug4 log was emitted below verbosity 4: %s", out)
	}
	if strings.Contains(out, "hidden debug5") {
		t.Fatalf("debug5 log was emitted below verbosity 5: %s", out)
	}

	buf.Reset()
	if err := SetDebugVerbosity(5); err != nil {
		t.Fatal(err)
	}
	Debug4("visible debug4")
	Debug5("visible debug5")
	out = buf.String()
	if !strings.Contains(out, "visible debug4") {
		t.Fatalf("debug4 log was not emitted at verbosity 5: %s", out)
	}
	if !strings.Contains(out, "visible debug5") {
		t.Fatalf("debug5 log was not emitted at verbosity 5: %s", out)
	}
}

func resetForTest() {
	defaultRuntime = newRuntimeState()
}
