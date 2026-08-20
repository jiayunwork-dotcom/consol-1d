package main

import (
	"fmt"
	"os"

	"consol-1d/internal/cli"
)

// main wires the two subcommands to the cli package. It holds no
// physics and no parsing: the consolidation kernel lives in internal/.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, cli.Usage())
		os.Exit(2)
	}
	switch os.Args[1] {
	case "profile":
		if err := cli.RunProfile(os.Args[2:], os.Stdout, os.Stderr); err != nil {
			cli.WriteError(os.Stderr, err)
			os.Exit(1)
		}
	case "curve":
		if err := cli.RunCurve(os.Args[2:], os.Stdout, os.Stderr); err != nil {
			cli.WriteError(os.Stderr, err)
			os.Exit(1)
		}
	case "settle":
		if err := cli.RunSettle(os.Args[2:], os.Stdout, os.Stderr); err != nil {
			cli.WriteError(os.Stderr, err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(cli.Usage())
	default:
		fmt.Fprintf(os.Stderr, "consol-1d: unknown command %q\n\n%s", os.Args[1], cli.Usage())
		os.Exit(2)
	}
}
