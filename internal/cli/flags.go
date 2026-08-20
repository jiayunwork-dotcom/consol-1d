package cli

import (
	"flag"
	"fmt"
	"strings"
)

// valueFlags describes the flags that take a value, so reorderFlags can
// move each flag together with its argument ahead of the positional
// ones. Go's flag package stops parsing at the first non-flag argument,
// but the documented invocation
// `consol-1d profile example/double-drain.json --t 1e6` puts flags after
// the scenario path, so the arguments are rearranged first.
var valueFlags = map[string]bool{
	"t":      true,
	"nodes":  true,
	"out":    true,
	"target": true,
	"times":  true,
}

// reorderFlags moves every flag (with its value) ahead of the positional
// arguments while preserving the relative order inside each group. A
// flag value is consumed greedily when the flag is known to take one
// and the next argument does not itself look like a flag.
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

// ProfileOptions is the parsed command line for the profile subcommand.
type ProfileOptions struct {
	ScenarioPath string
	TimeOverride *float64 // non-nil when --t was given explicitly
	Nodes        int
	OutFile      string
	JSONOutput   bool
}

// parseProfileFlags parses the remaining arguments after the "profile"
// subcommand word. Exactly one scenario path is required. The --t flag
// only overrides the JSON time when it is present on the command line,
// so a JSON-provided time cannot be silently clobbered by a zero
// default.
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
	// fs.Visit only walks flags that were set on the command line, which
	// is how an explicit --t (even --t 0) is told apart from the default.
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
