package cmdline

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gregoryv/cmdline/clitest"
)

func ExampleNewBasicParser() {
	os.Setenv("VERBOSE", "yes")
	var (
		cli      = NewBasicParser()
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password, $PASSWORD").String("")
		verbose  = cli.Flag("-V, --verbose, $VERBOSE")
		role     = cli.Option("-r, --role").Enum("guest", "admin", "nobody")
		url      = cli.Option("--test-host").Url("tcp://example.com:123")
		dur      = cli.Option("--pause").Duration("200ms")

		// parse and name non options
		username = cli.NamedArg("USERNAME").String("")
		note     = cli.NamedArg("NOTE").String("")
	)
	cli.Parse()

	// use options ...
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
	fmt.Fprintln(os.Stdout, uid, username, password, note, role, url, dur)
}

func ExampleNewParser() {
	var (
		cli      = NewParser()
		uid      = cli.Option("--uid", "Generated if not given").Int(0)
		password = cli.Option("-p, --password, $PASSWORD").String("")
		help     = cli.Flag("-h, --help")

		// parse and name non options
		username = cli.NamedArg("USERNAME").String("")
		note     = cli.NamedArg("NOTE").String("")
	)

	switch {
	case help:
		cli.Usage().WriteTo(os.Stdout)
		os.Exit(0)

	case !cli.Ok():
		fmt.Fprintln(os.Stderr, cli.Error())
		fmt.Fprintln(os.Stderr, "Try --help for more information")
		os.Exit(1)
	}

	// use options ...
	fmt.Fprintln(os.Stdout, uid, username, password, note)
}

func ExampleParser_Usage_optionalNamedArguments() {
	os.Args = []string{"mycmd"} // just for this test
	cli := NewBasicParser()
	cli.NamedArg("FILES...").Strings("file1", "file2")
	cli.Usage().WriteTo(os.Stdout)
	// output:
	//
	// Usage: mycmd [OPTIONS] [FILES...]
	//
	// Options
	//     -h, --help
}

func ExampleParser_Usage_requiredNamedArguments() {
	os.Args = []string{"mycmd", "-h"} // just for this test
	cli := NewBasicParser()
	cli.NamedArg("FILES...").Strings()
	cli.Usage().WriteTo(os.Stdout)
	// output:
	//
	// Usage: mycmd [OPTIONS] FILES...
	//
	// Options
	//     -h, --help
}
func Test_basic_parser_shows_help(t *testing.T) {
	cli := NewBasicParser()
	sh := clitest.NewShellT("test", "-h")
	cli.SetShell(sh)
	cli.Parse()
	if sh.ExitCode != 0 {
		t.Error(sh.Dump())
	}
	got := sh.Dump()
	if !strings.Contains(got, "Usage") {
		t.Error("dumped!")
	}
}

func Test_basic_parser_checks_errors(t *testing.T) {
	cli := NewBasicParser()
	sh := clitest.NewShellT("test", "-no-such")
	cli.SetShell(sh)
	log.SetOutput(ioutil.Discard)
	cli.Parse()
	if sh.ExitCode != 1 {
		t.Error(sh.Dump())
	}
}

func Test_parser_checks_errors(t *testing.T) {
	cli := NewParser()
	sh := clitest.NewShellT("test", "-no-such")
	cli.SetShell(sh)
	cli.Parse()
	if sh.ExitCode != 1 {
		t.Error(sh.Dump())
	}
}
func Test_parser_constructor_uses_osArgs(t *testing.T) {
	p := NewParser()
	if !reflect.DeepEqual(p.args, os.Args) {
		t.Fail()
	}
}

func Test_groups_are_unique(t *testing.T) {
	defer expectPanic(t)
	cli := Parse(t, "ls -h")
	cli.Group("actions", "selected")
	cli.Group("actions", "x")
}

func expectPanic(t *testing.T) {
	t.Helper()
	if e := recover(); e == nil {
		t.Error("should panic")
	}
}

func Test_no_arguments_is_ok(t *testing.T) {
	cli := Parse(t, "ls")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
	}
}

func Test_missing_last_value(t *testing.T) {
	cli := Parse(t, "mycmd -a=1 -b")
	cli.Option("-a").Int(0)
	cli.Option("-b").String("")
	if cli.Ok() {
		t.Fail()
	}
}

func Test_parser_string_option(t *testing.T) {
	cli := Parse(t, "mycmd -a=1 -b=k")
	cli.Option("-a").Int(0)
	got := cli.Option("-b").String("")
	if got != "k" {
		t.Fail()
	}
}

func Test_single_group_item_is_selected(t *testing.T) {
	cli := Parse(t, "mycmd hello")
	phrases := cli.Group("Phrases", "PHRASE")
	phrases.New("hello", nil)
	phrases.Selected()
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_unknown_group_item(t *testing.T) {
	cli := Parse(t, "mycmd car")
	nouns := cli.Group("Nouns", "NOUN")
	nouns.New("plane", nil)
	nouns.Selected()
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_required_multi_argument(t *testing.T) {
	cli := Parse(t, "touch")
	cli.NamedArg("FILES...").Strings()
	if cli.Ok() {
		t.Errorf("expected failure when required FILES... is missing")
	}
}

func Test_multi_argument_default_string(t *testing.T) {
	cli := Parse(t, "touch")
	got := cli.NamedArg("FILES...").String("file")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
	if got != "file" {
		t.Error("incorrect value:", got)
	}
}

func Test_Parser_reports_first_error(t *testing.T) {
	cli := Parse(t, "cmd -a=notint -b=1 -c=notint")
	cli.Option("-a").Int(0)
	cli.Option("-b").Int(0)
	cli.Option("-c").Int(0)
	err := cli.Error()
	if !strings.Contains(err.Error(), "-a") {
		t.Error(err)
	}
}

func Test_invalid_int_argument(t *testing.T) {
	cli := Parse(t, "cmd -i=k")
	cli.Option("-i").Int(10)
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_undefined_option(t *testing.T) {
	cli := Parse(t, "cmd -nosuch")
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_stringer(t *testing.T) {
	exp := "mycmd -help -i=4"
	cli := Parse(t, exp)
	got := cli.String()
	if !strings.Contains(got, exp) {
		t.Error("\n", exp, "\n", got)
	}
}

func TestParser_Argument_default_value(t *testing.T) {
	cli := Parse(t, "touch")
	got := cli.NamedArg("FILES...").Strings("file")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
	if len(got) == 0 || got[0] != "file" {
		t.Error("incorrect value:", got)
	}
}

func TestParser_Argument_multiple(t *testing.T) {
	cli := Parse(t, "touch a b")
	got := cli.NamedArg("FILES...").Strings()
	if !cli.Ok() {
		t.Error(cli.Error())
	}
	if got[1] != "b" {
		t.Error("incorrect value:", got)
	}
}

func TestParser_Argument_multiple_missing(t *testing.T) {
	cli := Parse(t, "touch")
	cli.NamedArg("FILES...").Strings()
	if cli.Ok() {
		t.Error("should fail")
	}
}
