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
		_, _ = cli.Group("Actions", &Hi{})
		_    = cli.Required("ACTION").String()
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
	}
	// output:
	// Usage: somecmd [OPTIONS] ACTION
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
	//
	// Actions
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

func (me *Hi) Run() error {
	fmt.Printf("Hi, %s!\n", me.to)
	return nil
}
