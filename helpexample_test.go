package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func Example_help() {
	var (
		cli  = cmdline.Parse("speak -h")
		_    = cli.Flag("-n, --dry-run")
		help = cli.Flag("-h, --help")

		// Group items
		phrases = cli.Group("Phrases", "PHRASE")

		// No extra options needed
		_ = phrases.New("askName", &Ask{})

		// Using builder function
		_ = phrases.New("sayHi", func(p *cmdline.Parser) interface{} {
			return &Hi{
				to: p.Option("-t, --to").String("stranger"),
			}
		})

		// Implementing the WithExtraOptions interface
		_ = phrases.New("compliment", &Compliment{})

		phrase = phrases.Selected()
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
	//
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
	fmt.Printf("%s, you look dashing I must say.", me.someone)
}

// ----------------------------------------

type runnable interface {
	Run()
}
