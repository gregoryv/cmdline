package main

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func main() {
	cli := cmdline.New(os.Args...)
	dryrun := cli.Flag("-n, --dry-run")
	name, opt := cli.Option("-firstname, --firstname").StringOpt("World")
	opt.Doc("The name of whom you want to say hello to")

	if err := cli.Error(); err != nil {
		if err != cmdline.ErrHelp {
			fmt.Println(err)
			os.Exit(1)
		}
		cli.WriteUsageTo(os.Stderr)
	}

	if !dryrun {
		fmt.Printf("Hello, %s!\n", name)
	}
}
