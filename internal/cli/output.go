package cli

import (
	"fmt"
	"io"
)

// WriteError prints a command failure to stderr in a stable form. The
// caller (main) then exits with a non-zero code; keeping the exit in
// main makes the wire-through trivial to test.
func WriteError(stderr io.Writer, err error) {
	fmt.Fprintf(stderr, "consol-1d: %v\n", err)
}

// Usage returns the top-level help text.
func Usage() string {
	return usage
}

// usage is the top-level help text shown without arguments or with
// "help". It documents the subcommand, the scenario fields and the
// failure semantics in one place.
const usage = `consol-1d: one-dimensional consolidation of a saturated soil layer.

Reads a Terzaghi consolidation scenario from a JSON file and reports the
pore-pressure dissipation profile, the average consolidation degree U
and the settlement ratio s/s_ult along the layer.

usage:
  consol-1d profile [-t seconds] [-nodes count] [-out file.csv] [-json] <scenario.json>
  consol-1d curve [-times a,b,c] [-json] <scenario.json>
  consol-1d settle [-target degree] <scenario.json>
  consol-1d help

scenario.json fields:
  cv              coefficient of consolidation, m^2/s (must be > 0)
  thickness       layer thickness H in metres (must be > 0)
  drainage        "single" (Hdr = H) or "double" (Hdr = H/2)
  initial_pressure {type: "uniform", u0} or {type: "linear", ua, ub}
  time            elapsed time t in seconds (must be >= 0)
  mv              optional volume compressibility, 1/kPa
  delta_sigma     optional stress increment, kPa

Tv = cv*t/Hdr^2 drives the Fourier series; U = 1 - mean(u)/u0 with the
same modes, and s = U*mv*delta_sigma*H when mv and delta_sigma are
given. Illegal parameters (cv <= 0, H <= 0, t < 0, negative initial
pressure, unknown drainage path) are reported on stderr and the process
exits non-zero.
`
