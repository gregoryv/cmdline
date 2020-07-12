package cmdline_test

import (
	"fmt"

	"github.com/gregoryv/cmdline"
)

func ExampleCommandLine() {
	args := []string{
		"adduser",
		"-n",
		"-p", "secret",
		"--uid", "100",
		"john",
	}

	cli := cmdline.New(args...)
	uid, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	password := cli.Option("-p, --password").String("")
	dryrun := cli.Flag("-n, --dry-run")
	username := cli.NeedArg("USERNAME", 0).String()
	cli.CheckOptions()

	fmt.Printf("uid=%v, username=%q, password=%q, dryrun=%v\n",
		uid, username, password, dryrun,
	)
	// output:
	// uid=100, username="john", password="secret", dryrun=true
}
