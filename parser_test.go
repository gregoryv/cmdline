package cmdline

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/golden"
)

func TestParser_AddGroup_twice(t *testing.T) {
	cli := Parse("ls -h")
	cli.Group("actions", "selected")
	defer catchPanic(t)
	cli.Group("actions", "x")
}

func catchPanic(t *testing.T) {
	t.Helper()
	e := recover()
	if e == nil {
		t.Error("should panic")
	}
}

func TestParser_Ok(t *testing.T) {
	cli := Parse("ls -r .")
	cli.Flag("-r")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
	}
}

func TestParser_not_Ok(t *testing.T) {
	args := "ls -r ."
	cli := Parse(args)
	cli.Option("-v").String("")
	phrases := cli.Group("Phrases", "PHRASE")
	phrases.New("hello", nil)
	phrases.Selected()
	if cli.Ok() {
		t.Log(cli.Error())
		t.Errorf("%q was ok, but -r is not defined", args)
	}
}

func TestParser_Required(t *testing.T) {
	cli := Parse("mkdir")
	cli.Required("DIR")
	if cli.Ok() {
		t.Errorf("expected failure when required DIR is missing")
	}
}

func TestParser_Optional(t *testing.T) {
	cli := Parse("ls")
	cli.Optional("DIR")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
	}
}
func TestParser_Usage(t *testing.T) {
	cli := NewParser("adduser")
	cli.Flag("-n, --dry-run")
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	cli.Option("-p, --password").String("")
	cli.Required("USERNAME").String()

	var buf bytes.Buffer
	cli.WriteUsageTo(&buf)
	golden.Assert(t, buf.String())
}

func TestParser_New_panic(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Error("Should panic on empty args")
		}
	}()
	NewParser()
}

func TestParser_Error(t *testing.T) {
	_, bad := asserter.NewErrors(t)

	bad(Parse("mycmd -h").Error())
	bad(Parse("mycmd -nosuch").Error())

	cli := Parse("cmd -i=k")
	cli.Option("-i").Int(10)
	bad(cli.Error())
}

func Test_Parser_reports_first_error(t *testing.T) {
	cli := Parse("cmd -a=notint -b=1 -c=notint")
	cli.Option("-a").Int(0)
	cli.Option("-b").Int(0)
	cli.Option("-c").Int(0)
	err := cli.Error()
	if !strings.Contains(err.Error(), "-a") {
		t.Error(err)
	}
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

func TestParser_Arg(t *testing.T) {
	cli := Parse("cp -i 1 /etc")
	cli.Option("-i").Int(0)
	assert := asserter.New(t)
	arg1 := cli.Required("FROM").String()
	assert().Equals(arg1, "/etc")
	arg2 := cli.Required("TO").String()
	assert().Equals(arg2, "")
}

func Test_stringer(t *testing.T) {
	cli := Parse("mycmd -help -i=4")
	assert := asserter.New(t)
	assert().Contains(cli.String(), "mycmd -help -i=4")
}
