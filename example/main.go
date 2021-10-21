package main

import (
	"fmt"

	"github.com/gregoryv/cmdline"
)

func main() {
	// Use the basic parser to automatically get a -h and --help flag
	cli := cmdline.NewBasicParser()
	cli.Preface(
		"speak - talks back to you",
		"Author: Gregory Vincic",
	)

	// Using a var block helps distinguishing setup from further execution
	var (
		// sub commands are common and are defined as a group
		phrases = cli.Group("Phrases", "PHRASE")

		// Bare sub command with no extra options
		_ = phrases.New("askName", runfunc(askName))

		// Using builder function that needs extra options
		_ = phrases.New("sayHi", func(p *cmdline.Parser) interface{} {
			return &Hi{
				to: p.Option("-t, --to").String("stranger"),
			}
		})

		// or you can implement the WithExtraOptions interface
		_ = phrases.New("compliment", &Compliment{})

		// select a phrase, if none is given it defaults to the first
		// one
		phrase = phrases.Selected()
	)

	// Use examples for common use cases. Note that you have to define
	// them before calling Parse.
	u := cli.Usage()
	u.Example(
		"Greet",
		"    $ speek sayHi -t John",
		"    Hi, John!",
	)
	cli.Parse()

	phrase.(runnable).Run()
}

func askName() {
	fmt.Println("What is your name?")
}

type Hi struct {
	to string
}

func (me *Hi) Run() { fmt.Printf("Hi, %s!\n", me.to) }

type Compliment struct {
	// Enable this subcommand to have a name
	cmdline.Item
	someone string
}

func (me *Compliment) ExtraOptions(p *cmdline.Parser) {
	me.someone = p.Option("-s, --someone").String("John")
}

func (me *Compliment) Run() {
	fmt.Printf("%s, you look dashing I must say.\n", me.someone)
}

// It is up to you to define the interface for your sub commands
type runnable interface {
	Run()
}

type runfunc func()

func (me runfunc) Run() { me() }
