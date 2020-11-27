package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example_help() {
	var (
		cli  = cmdline.Parse("speak sayHi -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")
	)
	// Add objects that implement Named and optional
	// WithExtraOptions interface
	phrases := cli.Group("Phrases", "PHRASE")
	phrases.New("askName", &Ask{})
	phrases.New("sayHi", func(cli *cmdline.Parser) interface{} {
		return &Hi{
			to: cli.Option("-t, --to").String("stranger"),
		}
	})
	phrases.New("compliment", &Compliment{})
	phrase := phrases.Selected()

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
	//     compliment
	//         -s, --someone : "John"

}

// ----------------------------------------

type Hi struct {
	to string
}

func (me *Hi) Run() { fmt.Printf("Hi, %s!\n", me.to) }

// ----------------------------------------

type Ask struct{ cmdline.Item }

func (me *Ask) Run() { fmt.Println("What is your name?") }

// ----------------------------------------

type Compliment struct {
	cmdline.Item
	someone string
}

// ExtraOptions
func (me *Compliment) ExtraOptions(p *cmdline.Parser) {
	me.someone = p.Option("-s, --someone").String("John")
}

func (me *Compliment) Run() {
	fmt.Printf("%s, don't you look nice today", me.someone)
}

// ----------------------------------------

type runnable interface {
	Run()
}
