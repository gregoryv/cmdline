package cmdline_test

import (
	"fmt"

	"xwing.7de.se/cmdline"
)

func ExampleCommandLine() {
	args := []string{
		"mycmd",
		"--username", "user",
		"-p", "secret",
		"--uid", "100",
		"-n",
	}

	cli := cmdline.New(args...)
	uid, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	username := cli.Option("-u, --username").String("john")
	password := cli.Option("-p, --password").String("")
	dryrun := cli.Flag("-n, --dry-run")
	cli.CheckOptions()

	fmt.Printf("uid=%v, username=%q, password=%q, dryrun=%v\n",
		uid, username, password, dryrun,
	)
	// output:
	// uid=100, username="user", password="secret", dryrun=true
}
