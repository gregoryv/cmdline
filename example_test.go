package cmdline_test

import (
	"fmt"
	"os"

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
	cli.Error()

	fmt.Printf("uid=%v, username=%q, password=%q, dryrun=%v\n",
		uid, username, password, dryrun,
	)
	// output:
	// uid=100, username="john", password="secret", dryrun=true
}

func Example() {
	cli := cmdline.Parse("somecmd -h")
	_ = cli.Option("-p, --password").String("")
	cli.Option("-i, --integer").IntOpt(0)
	_ = cli.NeedArg("HOST", 0).String()
	_ = cli.Flag("-n, --dry-run")

	help := cli.Flag("-h, --help")
	if err := cli.Error(); err != nil {
		fmt.Println(err)
		return
	}
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS] HOST
	//
	// Options
	//     -p, --password : ""
	//     -n, --dry-run : false
	//     -h, --help : false
}
