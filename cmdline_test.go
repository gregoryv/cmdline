package cmdline

import (
	"bytes"
	"reflect"
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

func TestCommandLine_Error(t *testing.T) {
	ok, bad := asserter.NewErrors(t)
	ok(Parse("cmd").Error())
	bad(Parse("mycmd -h").Error())
	bad(Parse("mycmd -nosuch").Error())

	cli := Parse("cmd -i=k")
	cli.Option("-i").Int(10)
	bad(cli.Error())
}

func Test_flags(t *testing.T) {
	args := "x -h a b"
	cli := Parse(args)
	if !cli.Flag("-h") {
		t.Errorf("-h failed for %q", args)
	}
	if cli.Flag("--h") {
		t.Errorf("--h was ok for %q", args)
	}
	got := cli.Args()
	if !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Error(got)
	}
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
