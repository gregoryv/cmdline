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
		cli = cmdline.New(args...)
		uid = cli.Option("--uid",
			"user id to set on the new account",
			"If not given, one is generated",
		).Int(0)
		password = cli.Option("-p, --password").String("")
		help     = cli.Flag("-h, --help") // explicit help handling

		// parse and name non options
		username = cli.NeedArg("USERNAME").String()
	)

	switch {
	case !cli.Ok():
		fmt.Println(cli.Error())
		fmt.Println("Try --help for more information")

	case help:
		cli.WriteUsageTo(os.Stdout)

	default:
		fmt.Println(uid, username, password)
	}
}

func Example_help() {
	var (
		cli  = cmdline.Parse("somecmd -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS]
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
}
