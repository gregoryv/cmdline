package cmdline_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
)

func ExampleUsage() {
	cli := cmdline.NewBasicParser()
	cli.SetArgs("speak")
	var (
		_ = cli.Flag("-n, --dry-run")
		_ = cli.Option("-u, --username, $USER").String("")

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
	)
	cli.Usage().WriteTo(os.Stdout)

	// output:
	// Usage: speak [OPTIONS] [PHRASE]
	//
	// Options
	//     -n, --dry-run
	//     -u, --username, $USER : ""
	//     -h, --help
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
