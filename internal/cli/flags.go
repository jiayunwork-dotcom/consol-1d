package cli

import (
	"flag"
	"fmt"
	"strings"
)

var valueFlags = map[string]bool{
	"t":      true,
	"nodes":  true,
	"out":    true,
	"target": true,
	"times":  true,
}

func reorderFlags(args []string) []string {
	var flags, positional []string
	for i := 0; i < len(args); {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
			i++
			continue
		}
		name := strings.TrimLeft(a, "-")
		if eq := strings.Index(name, "="); eq >= 0 {
			flags = append(flags, a)
			i++
			continue
		}
		flags = append(flags, a)
		if valueFlags[name] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			flags = append(flags, args[i+1])
			i++
		}
		i++
	}
	return append(flags, positional...)
}

type ProfileOptions struct {
	ScenarioPath string
	TimeOverride *float64
	Nodes        int
	OutFile      string
	JSONOutput   bool
}

func parseProfileFlags(args []string) (ProfileOptions, error) {
	args = reorderFlags(args)
	fs := flag.NewFlagSet("profile", flag.ContinueOnError)
	t := fs.Float64("t", 0, "override the scenario time in seconds (default: keep the JSON time)")
	nodes := fs.Int("nodes", 11, "number of profile nodes, at least 2")
	out := fs.String("out", "", "write the profile to this CSV file")
	jsonOut := fs.Bool("json", false, "print the result as JSON instead of the human report")
	if err := fs.Parse(args); err != nil {
		return ProfileOptions{}, err
	}
	if fs.NArg() != 1 {
		return ProfileOptions{}, fmt.Errorf("profile needs exactly one scenario JSON file, got %d arguments", fs.NArg())
	}
	opt := ProfileOptions{
		ScenarioPath: fs.Arg(0),
		Nodes:        *nodes,
		OutFile:      *out,
		JSONOutput:   *jsonOut,
	}
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if set["t"] {
		opt.TimeOverride = t
	}
	if opt.Nodes < 2 {
		return ProfileOptions{}, fmt.Errorf("profile needs at least 2 nodes, got %d", opt.Nodes)
	}
	return opt, nil
}
