package report_test

import (
	"strings"
	"testing"

	"consol-1d/internal/model"
	"consol-1d/internal/report"
)

// TestReportPrintsKeyNumbers checks that the human report carries the
// three headline numbers: U, the midpoint pressure and the settlement
// ratio.
func TestReportPrintsKeyNumbers(t *testing.T) {
	res, err := model.Solve(model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7), 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	text := report.Report(res)
	for _, want := range []string{"average U", "midpoint", "settlement", "profile"} {
		if !strings.Contains(text, want) {
			t.Errorf("report should contain %q, got:\n%s", want, text)
		}
	}
}

// TestProfileTableRowCount checks that the table renders the header
// plus exactly one data row per profile point.
func TestProfileTableRowCount(t *testing.T) {
	res, err := model.Solve(model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7), 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	table := report.ProfileTable(res.Profile)
	rows := 0
	for _, line := range strings.Split(table, "\n") {
		if strings.TrimSpace(line) != "" {
			rows++
		}
	}
	if want := len(res.Profile) + 1; rows != want {
		t.Errorf("profile table rows (header + points): expected %d, got %d", want, rows)
	}
}

// TestCSVHasHeaderAndRows checks the CSV dump starts with the header
// and carries every point as a data row.
func TestCSVHasHeaderAndRows(t *testing.T) {
	res, err := model.Solve(model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7), 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	csv := report.CSV(res.Profile)
	lines := strings.Split(strings.TrimRight(csv, "\n"), "\n")
	if want := "z_over_H,z_m,u_kPa,u_over_u0,dissipated"; lines[0] != want {
		t.Errorf("csv header: expected %q, got %q", want, lines[0])
	}
	if want := len(res.Profile); len(lines)-1 != want {
		t.Errorf("csv data rows: expected %d, got %d", want, len(lines)-1)
	}
}

// TestJSONDocumentContainsHeadlineNumbers checks the structured output
// round-trips through the JSON decoder with the key scalars intact.
func TestJSONDocumentContainsHeadlineNumbers(t *testing.T) {
	res, err := model.Solve(model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7), 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	data, err := report.ToJSON(res)
	if err != nil {
		t.Fatalf("ToJSON: expected no error, got %v", err)
	}
	s := string(data)
	if !strings.Contains(s, `"average_consolidation"`) || !strings.Contains(s, `"midpoint_pressure"`) {
		t.Errorf("json should contain the headline fields, got:\n%s", s)
	}
	if !strings.Contains(s, `"drainage_distance"`) {
		t.Errorf("json should expose the drainage distance, got:\n%s", s)
	}
}
