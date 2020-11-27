package cmdline_test

import (
	"fmt"

	"github.com/gregoryv/cmdline"
)

/*
This example shows how to group a set of options for a sub command.
*/
func Example_groupedSubCommands() {
	var (
		cli    = cmdline.Parse("speach sayHi --to Gopher")
		phrase = cli.Group("Phrases", "PHRASE", &Hi{}).Item()
	)

	if !cli.Ok() {
		fmt.Println(cli.Error())
		return
	}

	phrase.ExtraOptions(cli)
	phrase.(Runnable).Run()
	// output:
	// Hi, Gopher!
}

type Runnable interface {
	Run() error
}
