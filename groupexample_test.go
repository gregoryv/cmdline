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
		phrases, _ = cli.Group("Phrases", &Hi{})
		name       = cli.Required("PHRASE").String()
	)
	phrase, found := phrases.Find(name)
	if !found {
		fmt.Println("no such phrase")
		os.Exit(1)
	}

	phrase.ExtraOptions(cli)
	phrase.(Runnable).Run()
	// output:
	// Hi, Gopher!
}

type Runnable interface {
	Run() error
}
