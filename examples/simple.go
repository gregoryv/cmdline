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
	cli.CheckOptions()

	if !dryrun {
		fmt.Println("Hello", name)
	}
}
