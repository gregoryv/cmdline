package cmdline

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/golden"
)

func TestCommandLine_Usage(t *testing.T) {
	cli := New("adduser")
	cli.Flag("-n, --dry-run")
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	cli.Option("-p, --password").String("")
	cli.NeedArg("USERNAME", 0).String()

	var buf bytes.Buffer
	cli.WriteUsageTo(&buf)
	golden.Assert(t, buf.String())
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
	assert().Equals(cli.Argn(0), "a")
	assert().Equals(cli.Argn(3), "")
}

func TestCommandLine_Arg(t *testing.T) {
	cli := Parse("cp -i 1 /etc")
	cli.Option("-i").Int(0)
	assert := asserter.New(t)
	arg1 := cli.NeedArg("FROM", 0).String()
	assert().Equals(arg1, "/etc")
	arg2 := cli.NeedArg("TO", 1).String()
	assert().Equals(arg2, "")
}

func Test_stringer(t *testing.T) {
	cli := Parse("mycmd -help -i=4")
	assert := asserter.New(t)
	assert().Contains(cli.String(), "mycmd -help -i=4")
}
