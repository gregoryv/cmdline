package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example() {
	// would be os.Args...
	args := []string{"adduser", "-h", "-p", "secret", "--uid", "100", "john"}

	var (
		cli      = cmdline.New(args...)
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.Required("USERNAME").String()
		note     = cli.Optional("NOTE").String()
	)

	switch {
	case !cli.Ok():
		fmt.Println(cli.Error())
		fmt.Println("Try --help for more information")

	case help:
		cli.WriteUsageTo(os.Stdout)

	default:
		fmt.Println(uid, username, password, note)
	}
}

func Example_help() {
	var (
		cli  = cmdline.Parse("somecmd -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")
		// order is important for non options
		_ = cli.Required("FILE")
		_ = cli.Optional("DIR")
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS] FILE [DIR]
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
}
