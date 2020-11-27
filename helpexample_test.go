package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example_help() {
	var (
		cli    = cmdline.Parse("speak sayHi -h")
		_      = cli.Flag("-n, --dry-run")
		help   = cli.Flag("-h, --help")
		phrase = cli.Group("Phrases", "PHRASE",
			// Add objects that implement Named and optional
			// WithExtraOptions interface

			&Ask{ // first is default
				Item: "askName",
			},
			&Hi{
				Item: "sayHi",
				to:   "stranger",
			},
		).Item()
	)
	if help {
		cli.WriteUsageTo(os.Stdout)
		return
	}
	if !cli.Ok() {
		fmt.Println(cli.Error())
		return
	}

	phrase.(runnable).Run()
	// output:
	// Usage: speak [OPTIONS] [PHRASE]
	//
	// Options
	//     -n, --dry-run : false
	//     -h, --help : false
	//
	// Phrases
	//     askName (default)
	//     sayHi
	//         -t, --to : "stranger"
}

// ----------------------------------------

type Hi struct {
	cmdline.Item // implements the Named interface
	to           string
}

// ExtraOptions implements WithExtraOptions interface
func (me *Hi) ExtraOptions(cli *cmdline.Parser) {
	me.to = cli.Option("-t, --to").String(me.to)
}

func (me *Hi) Run() { fmt.Printf("Hi, %s!\n", me.to) }

// ----------------------------------------

type Ask struct{ cmdline.Item }

func (me *Ask) Run() { fmt.Printf("What is your name?") }

// ----------------------------------------

type runnable interface {
	Run()
}
