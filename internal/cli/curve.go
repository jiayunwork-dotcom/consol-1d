package cli

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"consol-1d/internal/model"
	"consol-1d/internal/report"
)

// RunCurve executes the curve subcommand: it evaluates the scenario at
// a batch of times and prints the consolidation curve. The default time
// list spans the consolidation process; `-times` accepts a
// comma-separated list of seconds.
func RunCurve(args []string, stdout, stderr io.Writer) error {
	args = reorderFlags(args)
	fs := flag.NewFlagSet("curve", flag.ContinueOnError)
	timesList := fs.String("times", "", "comma-separated times in seconds (default: auto-spanned)")
	nodes := fs.Int("nodes", 2, "profile nodes used per solve")
	jsonOut := fs.Bool("json", false, "print each curve point as JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("curve needs exactly one scenario JSON file, got %d arguments", fs.NArg())
	}
	in, err := LoadInput(fs.Arg(0))
	if err != nil {
		return err
	}
	var times []float64
	if *timesList != "" {
		times, err = parseTimes(*timesList)
		if err != nil {
			return err
		}
	} else {
		times = model.DefaultCurveTimes(in)
	}
	pts, err := model.ConsolidationCurve(in, times, *nodes)
	if err != nil {
		return err
	}
	if *jsonOut {
		for _, p := range pts {
			if _, err := fmt.Fprintf(stdout,
				`{"t":%g,"Tv":%g,"U":%g,"u_mid_kPa":%g,"s_m":%g}`+"\n",
				p.Time, p.TimeFactor, p.U, p.MidpointPressure, p.Settlement); err != nil {
				return err
			}
		}
		return nil
	}
	if _, err := fmt.Fprint(stdout, report.CurveTable(pts)); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, report.CurveSummary(pts))
	return err
}

// parseTimes splits and parses a comma-separated list of seconds.
func parseTimes(list string) ([]float64, error) {
	parts := strings.Split(list, ",")
	times := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("bad time %q: %w", p, err)
		}
		times = append(times, v)
	}
	if len(times) == 0 {
		return nil, fmt.Errorf("no usable times in %q", list)
	}
	return times, nil
}

// RunSettle executes the settle subcommand: it reports the time at
// which the layer reaches a target average consolidation degree.
func RunSettle(args []string, stdout, stderr io.Writer) error {
	args = reorderFlags(args)
	fs := flag.NewFlagSet("settle", flag.ContinueOnError)
	target := fs.Float64("target", 0.9, "target average consolidation degree in (0, 1)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("settle needs exactly one scenario JSON file, got %d arguments", fs.NArg())
	}
	in, err := LoadInput(fs.Arg(0))
	if err != nil {
		return err
	}
	t, err := model.TimeToDegree(in, *target, 1e-9)
	if err != nil {
		return err
	}
	path := model.NewDrainagePath(in.Drainage, in.Thickness)
	tv := t * in.Cv / (path.Hdr * path.Hdr)
	_, err = fmt.Fprintf(stdout,
		"settle: time to U=%g is %s s (%s); Tv=%s, Hdr=%g m\n",
		*target, report.Sci(t), report.Duration(t), report.Sci(tv), path.Hdr)
	return err
}
