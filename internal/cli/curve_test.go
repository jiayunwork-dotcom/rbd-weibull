package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCurveCommand(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunCurve([]string{path, "--from", "0", "--to", "4500", "--n", "4"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("RunCurve exit = %d, stderr: %s", code, errBuf.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("RunCurve printed %d lines, want 4:\n%s", len(lines), out.String())
	}
	if !strings.Contains(lines[0], "0 ") || !strings.Contains(lines[3], "4500") {
		t.Fatalf("curve endpoints wrong:\n%s", out.String())
	}
}

func TestCurveRejectsBadHorizon(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunCurve([]string{path, "--to", "-5"}, &out, &errBuf)
	if code != 2 {
		t.Fatalf("RunCurve exit = %d, want 2", code)
	}
}

func TestDescribeCommand(t *testing.T) {
	path := writeTemp(t, "pump2.json", pump2JSON)
	var out, errBuf bytes.Buffer
	code := RunDescribe([]string{path}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("RunDescribe exit = %d, stderr: %s", code, errBuf.String())
	}
	text := out.String()
	for _, want := range []string{"series (2 blocks)", "parallel (2 blocks)", "pump-A", "valve"} {
		if !strings.Contains(text, want) {
			t.Fatalf("describe output missing %q:\n%s", want, text)
		}
	}
}
