package cmdline_test

import (
	"fmt"
	"io"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example() {
	run(os.Stdout, "adduser", "-p", "secret", "--uid", "100", "john")
}

func run(w io.Writer, args ...string) {
	var (
		cli      = cmdline.NewParser(args...)
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.Required("USERNAME").String()
		note     = cli.Optional("NOTE").String()
	)

	switch {
	case help:
		cli.WriteUsageTo(w)

	case !cli.Ok():
		fmt.Fprintln(w, cli.Error())
		fmt.Fprintln(w, "Try --help for more information")

	default:
		fmt.Fprintln(w, uid, username, password, note)
	}
}
