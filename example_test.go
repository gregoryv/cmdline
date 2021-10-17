package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example() {
	var (
		cli      = cmdline.NewParser()
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password, $PASSWORD").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.Required("USERNAME").String("")
		note     = cli.Optional("NOTE").String("")
	)

	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)
		os.Exit(0)

	case !cli.Ok():
		fmt.Fprintln(os.Stderr, cli.Error())
		fmt.Fprintln(os.Stderr, "Try --help for more information")
		os.Exit(1)
	}

	// use options ...
	fmt.Fprintln(os.Stdout, uid, username, password, note)
}
