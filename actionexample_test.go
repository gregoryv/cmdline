package cmdline_test

import "github.com/gregoryv/cmdline"

/*
This example shows how to group a set of options for a sub command.
*/
func Example_groupedSubCommands() {
	var (
		cli        = cmdline.Parse("somecmd sayHi --to Gopher")
		actions, _ = cli.Group("Actions", &Hi{})
		name       = cli.Required("ACTION").String()
	)
	action := actions.FindAction(name)
	action.ExtraOptions(cli)
	action.(Runnable).Run()
	// output:
	// Hi, Gopher!
}

type Runnable interface {
	Run() error
}
