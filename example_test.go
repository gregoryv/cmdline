package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func ExampleNewBasicParser() {
	var (
		cli      = cmdline.NewBasicParser()
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password, $PASSWORD").String("")

		// parse and name non options
		username = cli.Required("USERNAME").String("")
		note     = cli.Optional("NOTE").String("")
	)
	cli.Parse()

	// use options ...
	fmt.Fprintln(os.Stdout, uid, username, password, note)
}

func ExampleNewParser() {
	var (
		cli      = cmdline.NewParser()
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password, $PASSWORD").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.Argument("USERNAME").String("")
		note     = cli.Optional("NOTE").String("")
	)

	switch {
	case help:
		cli.Usage().WriteTo(os.Stdout)
		os.Exit(0)

	case !cli.Ok():
		fmt.Fprintln(os.Stderr, cli.Error())
		fmt.Fprintln(os.Stderr, "Try --help for more information")
		os.Exit(1)
	}

	// use options ...
	fmt.Fprintln(os.Stdout, uid, username, password, note)
}

func ExampleParser_Usage() {
	os.Args = []string{"mycmd"} // just for this test
	cli := cmdline.NewBasicParser()
	cli.Argument("FILES...").Strings("file1", "file2")
	cli.Usage().WriteTo(os.Stdout)
	// output:
	//
	// Usage: mycmd [OPTIONS] [FILES...]
	//
	// Options
	//     -h, --help
}
