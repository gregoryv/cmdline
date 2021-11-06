package cmdline

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_basic_parser_shows_help(t *testing.T) {
	cli := NewBasicParser()
	sh := NewShellT("test", "-h")
	cli.SetShell(sh)
	cli.Parse()
	if sh.ExitCode != 0 {
		t.Error(sh.Dump())
	}
}

func Test_basic_parser_checks_errors(t *testing.T) {
	cli := NewBasicParser()
	sh := NewShellT("test", "-no-such")
	cli.SetShell(sh)
	log.SetOutput(ioutil.Discard)
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

func Test_required_argument(t *testing.T) {
	cli := Parse(t, "mkdir")
	cli.Required("DIR")
	if cli.Ok() {
		t.Errorf("expected failure when required DIR is missing")
	}
}

func Test_optional_argument(t *testing.T) {
	cli := Parse(t, "ls")
	cli.Optional("DIR")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
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
