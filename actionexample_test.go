package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

/*
This example shows how to group a set of options for a sub command.
*/
func Example_groupedSubCommands() {
	var (
		cli        = cmdline.Parse("speach sayHi --to Gopher")
		actions, _ = cli.Group("Actions", &Hi{})
		name       = cli.Required("ACTION").String()
	)
	action, found := actions.Find(name)
	if !found {
		fmt.Println("no such action")
		os.Exit(1)
	}

	action.ExtraOptions(cli)
	action.(Runnable).Run()
	// output:
	// Hi, Gopher!
}

type Runnable interface {
	Run() error
}
