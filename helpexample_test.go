package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example_help() {
	var (
		cli  = cmdline.Parse("somecmd -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")
		// order is important for non options
		_ = cli.Group("Actions", "ACTION",
			&Ask{}, // first is always default
			&Hi{},
		)
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS] [ACTION]
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
	//
	// Actions
	//     askName (default)
	//     sayHi
	//         -t, --to : "stranger"
}

// Hi implements the sayHi action
type Hi struct {
	to string
}

func (me *Hi) Name() string { return "sayHi" }

func (me *Hi) ExtraOptions(cli *cmdline.Parser) {
	me.to = cli.Option("-t, --to").String("stranger")
}

func (me *Hi) Run() { fmt.Printf("Hi, %s!\n", me.to) }

type Ask struct{}

func (me *Ask) Name() string { return "askName" }

func (me *Ask) Run() { fmt.Printf("What is your name?") }
