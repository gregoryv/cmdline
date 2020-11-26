package cmdline

import (
	"os"
)

func Example_help() {
	var (
		cli  = Parse("somecmd -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")
		// order is important for non options
		_ = cli.Required("ACTION")
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS] ACTION
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
}
