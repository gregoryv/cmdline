package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example() {
	var (
		cli = cmdline.NewParser(
			// os.Args...
			"adduser", "-p", "secret", "--uid", "100", "john",
		)
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.Required("USERNAME").String("")
		note     = cli.Optional("NOTE").String("")
	)

	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)

	case !cli.Ok():
		fmt.Fprintln(os.Stderr, cli.Error())
		fmt.Fprintln(os.Stderr, "Try --help for more information")

	default:
		fmt.Fprintln(os.Stdout, uid, username, password, note)
	}
}
