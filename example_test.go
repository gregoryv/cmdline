package cmdline_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/cmdline/clitest"
	"github.com/gregoryv/golden"
)

func TestUsage_WriteTo(t *testing.T) {
	cli := cmdline.NewBasicParser()

	sh := clitest.NewShellT("speak", "-h")
	cli.SetShell(sh)

	cli.Preface(
		"speak - talks back to you",
		"Author: Gregory Vincic",
	)
	var (
		_ = cli.Flag("-n, --dry-run")
		_ = cli.Option("-u, --username, $USER").String("")
		_ = cli.Option("-r, --role").Enum("user", "user", "admin")

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
	u := cli.Usage()
	u.Example(
		"Greet",
		"    $ speek sayHi -t John",
		"    Hi, John!",
	)
	cli.Parse()
	var buf bytes.Buffer
	cli.Usage().WriteTo(&buf)
	golden.Assert(t, buf.String())
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
