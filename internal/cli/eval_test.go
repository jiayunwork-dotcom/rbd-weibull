package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

const pump2JSON = `{
  "type": "series",
  "blocks": [
    { "type": "parallel", "blocks": [
      { "type": "unit", "name": "pump-A", "beta": 1.8, "eta": 4500 },
      { "type": "unit", "name": "pump-B", "beta": 1.8, "eta": 4500 }
    ]},
    { "type": "unit", "name": "valve", "beta": 2.5, "eta": 20000 }
  ]
}`

func TestEvalCommand(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunEval([]string{path, "--t", "1000"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("RunEval exit = %d, stderr: %s", code, errBuf.String())
	}
	text := out.String()
	for _, want := range []string{"system R(t)", "system MTTF", "pump-A", "valve"} {
		if !strings.Contains(text, want) {
			t.Fatalf("output missing %q:\n%s", want, text)
		}
	}
	// The pump2 diagram at t=1000 gives R ≈ 0.995277.
	if !strings.Contains(text, "0.99527") {
		t.Fatalf("output does not carry the expected reliability:\n%s", text)
	}
}

func TestEvalMissionTimeFlag(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunEval([]string{"--t=500", path}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("RunEval exit = %d, stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "500") {
		t.Fatalf("output does not echo the mission time:\n%s", out.String())
	}
}

func TestEvalBadFile(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunEval([]string{"no-such-file.json"}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("RunEval exit = %d, want 1", code)
	}
	if errBuf.Len() == 0 {
		t.Fatal("stderr is empty for a missing input file")
	}
}

func TestEvalValidationError(t *testing.T) {
	bad := `{ "type": "kofn", "k": 3, "blocks": [
	  { "type": "unit", "beta": 1, "eta": 10 },
	  { "type": "unit", "beta": 1, "eta": 10 }
	] }`
	path := writeTemp(t, "bad.json", bad)
	var out, errBuf bytes.Buffer
	code := RunEval([]string{path}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("RunEval exit = %d, want 1", code)
	}
	if !strings.Contains(errBuf.String(), "k greater than number of blocks") {
		t.Fatalf("stderr does not mention the k>n violation: %s", errBuf.String())
	}
}

func TestEvalZeroTimeRejected(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunEval([]string{path, "--t", "0"}, &out, &errBuf)
	if code != 2 {
		t.Fatalf("RunEval exit = %d, want 2 for a non positive time", code)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run([]string{"frobnicate"}, &out, &errBuf)
	if code != 2 {
		t.Fatalf("Run exit = %d, want 2", code)
	}
	if !strings.Contains(errBuf.String(), "unknown command") {
		t.Fatalf("stderr does not mention the unknown command: %s", errBuf.String())
	}
}
