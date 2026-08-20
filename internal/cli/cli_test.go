package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"consol-1d/internal/cli"
)

const validScenario = `{
  "cv": 1e-7,
  "thickness": 10.0,
  "drainage": "double",
  "initial_pressure": {"type": "uniform", "u0": 100.0},
  "time": 50000000.0,
  "mv": 0.0002,
  "delta_sigma": 100.0
}
`

func writeScenario(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "scenario.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write scenario: expected no error, got %v", err)
	}
	return path
}

// TestProfileRejectsBadJSON checks that malformed scenario JSON fails
// the subcommand instead of being half-parsed.
func TestProfileRejectsBadJSON(t *testing.T) {
	path := writeScenario(t, `{"cv": 1e-7, "thick`)
	var out, errBuf bytes.Buffer
	if err := cli.RunProfile([]string{path}, &out, &errBuf); err == nil {
		t.Error("RunProfile on broken JSON: expected an error, got nil")
	}
}

// TestProfileRejectsNegativeTime checks that a negative --t override is
// rejected just like a negative time inside the JSON.
func TestProfileRejectsNegativeTime(t *testing.T) {
	path := writeScenario(t, validScenario)
	var out, errBuf bytes.Buffer
	if err := cli.RunProfile([]string{path, "--t", "-1"}, &out, &errBuf); err == nil {
		t.Error("RunProfile with --t -1: expected an error, got nil")
	}
}

// TestProfileOverrideTimePrintsSummary checks the documented invocation
// form (scenario path before the flags) and the presence of the
// headline numbers in the output.
func TestProfileOverrideTimePrintsSummary(t *testing.T) {
	path := writeScenario(t, validScenario)
	var out, errBuf bytes.Buffer
	if err := cli.RunProfile([]string{path, "--t", "1e6"}, &out, &errBuf); err != nil {
		t.Fatalf("RunProfile: expected no error, got %v", err)
	}
	for _, want := range []string{"U=", "settlement_ratio", "midpoint"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output should contain %q, got:\n%s", want, out.String())
		}
	}
}

// TestProfileWritesCSV checks the -out flag persists the profile dump
// with its header line.
func TestProfileWritesCSV(t *testing.T) {
	path := writeScenario(t, validScenario)
	csvPath := filepath.Join(t.TempDir(), "profile.csv")
	var out, errBuf bytes.Buffer
	if err := cli.RunProfile([]string{"-out", csvPath, path}, &out, &errBuf); err != nil {
		t.Fatalf("RunProfile: expected no error, got %v", err)
	}
	data, err := os.ReadFile(csvPath)
	if err != nil {
		t.Fatalf("read csv: expected no error, got %v", err)
	}
	if !strings.HasPrefix(string(data), "# consol-1d:") {
		t.Errorf("csv should start with the preamble, got:\n%s", data)
	}
	if !strings.Contains(string(data), "z_over_H,z_m,u_kPa,u_over_u0,dissipated") {
		t.Errorf("csv should contain the header line, got:\n%s", data)
	}
}

// TestParseInputRejectsUnknownField checks the strict decoder refuses a
// misspelled scenario key instead of silently ignoring it.
func TestParseInputRejectsUnknownField(t *testing.T) {
	bad := `{"cv": 1e-7, "thickness": 10, "drainage": "double",
	  "initial_pressure": {"type": "uniform", "u0": 100},
	  "time": 5e7, "thicknes": 9}`
	if _, err := cli.ParseInput([]byte(bad)); err == nil {
		t.Error("ParseInput with unknown field: expected an error, got nil")
	}
}

// TestParseInputAcceptsMinimalScenario checks the minimal legal JSON is
// accepted and carries the expected time factor after a solve.
func TestParseInputAcceptsMinimalScenario(t *testing.T) {
	path := writeScenario(t, validScenario)
	in, err := cli.LoadInput(path)
	if err != nil {
		t.Fatalf("LoadInput: expected no error, got %v", err)
	}
	if in.Cv != 1e-7 || in.Thickness != 10 {
		t.Errorf("loaded scenario: expected cv=1e-7 and H=10, got %g and %g", in.Cv, in.Thickness)
	}
}

// TestCurvePrintsMonotoneConsolidation checks the curve subcommand
// reports a rising U over increasing times.
func TestCurvePrintsMonotoneConsolidation(t *testing.T) {
	path := writeScenario(t, validScenario)
	var out, errBuf bytes.Buffer
	if err := cli.RunCurve([]string{path, "-times", "1e5,1e6,1e7"}, &out, &errBuf); err != nil {
		t.Fatalf("RunCurve: expected no error, got %v", err)
	}
	if !strings.Contains(out.String(), "U") || !strings.Contains(out.String(), "curve:") {
		t.Errorf("curve output should carry a table and a summary, got:\n%s", out.String())
	}
}

// TestSettleReportsTimeToTarget checks the settle subcommand returns a
// finite time for a mid-range target degree and states the time factor.
func TestSettleReportsTimeToTarget(t *testing.T) {
	path := writeScenario(t, validScenario)
	var out, errBuf bytes.Buffer
	if err := cli.RunSettle([]string{path, "-target", "0.5"}, &out, &errBuf); err != nil {
		t.Fatalf("RunSettle: expected no error, got %v", err)
	}
	for _, want := range []string{"settle:", "Tv="} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("settle output should contain %q, got:\n%s", want, out.String())
		}
	}
}

// TestSettleRejectsOutOfRangeTarget checks a target outside (0, 1) is
// refused before any time is computed.
func TestSettleRejectsOutOfRangeTarget(t *testing.T) {
	path := writeScenario(t, validScenario)
	var out, errBuf bytes.Buffer
	if err := cli.RunSettle([]string{path, "-target", "1.5"}, &out, &errBuf); err == nil {
		t.Error("RunSettle with target 1.5: expected an error, got nil")
	}
}
