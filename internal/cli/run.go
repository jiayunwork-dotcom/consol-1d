package cli

import (
	"fmt"
	"io"
	"os"

	"consol-1d/internal/model"
	"consol-1d/internal/report"
)

func RunProfile(args []string, stdout, stderr io.Writer) error {
	opt, err := parseProfileFlags(args)
	if err != nil {
		return err
	}
	in, err := LoadInput(opt.ScenarioPath)
	if err != nil {
		return err
	}
	if opt.TimeOverride != nil {
		in.Time = *opt.TimeOverride
	}
	res, err := model.Solve(in, opt.Nodes)
	if err != nil {
		return err
	}
	if opt.OutFile != "" {
		if err := writeCSV(opt.OutFile, res); err != nil {
			return err
		}
	}
	if opt.JSONOutput {
		data, err := report.ToJSON(res)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, string(data))
		return err
	}
	_, err = fmt.Fprint(stdout, report.Report(res))
	if err == nil {
		_, err = fmt.Fprintln(stdout, report.Summary(res))
	}
	return err
}

func writeCSV(path string, res model.Result) error {
	data := HoldCSVLive(report.CSVOfResult(res))
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(data), 0o644); err != nil {
		return fmt.Errorf("write profile CSV %s: %w", path, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("finalize profile CSV %s: %w", path, err)
	}
	return nil
}
