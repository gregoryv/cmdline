package cmdline

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/golden"
)

func TestCommandLine_Usage(t *testing.T) {
	cli := New("mycmd")
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	cli.Option("-u, --username").String("john")
	cli.Option("-p, --password").String("")
	cli.Flag("-n, --dry-run")
	var buf bytes.Buffer
	cli.WriteUsageTo(&buf)
	got := buf.String()
	golden.Assert(t, got)
}

func TestCommandLine_New_panic(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("Should panic on empty args")
		}
	}()
	New()
}

func TestCommandLine_CheckOptions(t *testing.T) {
	ok, bad := asserter.NewErrors(t)
	ok(Parse("cmd").CheckOptions())
	bad(Parse("mycmd -h").CheckOptions())
	bad(Parse("mycmd -nosuch").CheckOptions())

	cli := Parse("cmd -i=k")
	cli.Option("-i").Int(10)
	bad(cli.CheckOptions())
}

func TestCommandLine_Args(t *testing.T) {
	cli := Parse("x -h -i=3 a b")
	cli.Option("-h").Bool()
	cli.Option("-i").Int(0)
	assert := asserter.New(t)
	assert().Equals(len(cli.Args()), 2)
	assert().Equals(cli.Arg(0), "a")
	assert().Equals(cli.Arg(3), "")
}

func Test_stringer(t *testing.T) {
	cli := Parse("mycmd -help -i=4")
	assert := asserter.New(t)
	assert().Contains(cli.String(), "mycmd -help -i=4")
}
