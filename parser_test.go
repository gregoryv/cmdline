package cmdline

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_parser_checks_errors(t *testing.T) {
	cli := NewParser()
	cli.SetArgs("test", "-no-such")
	var code int
	cli.SetExit(func(v int) { code = v })
	var (
		_ = cli.Flag("-h, --help")
	)
	cli.Parse()
	if code != 1 {
		t.Error("exit code:", code)
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
	cli := Parse("ls -h")
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
	cli := Parse("ls")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
	}
}

func Test_missing_last_value(t *testing.T) {
	cli := Parse("mycmd -a=1 -b")
	cli.Option("-a").Int(0)
	cli.Option("-b").String("")
	if cli.Ok() {
		t.Fail()
	}
}

func Test_parser_string_option(t *testing.T) {
	cli := Parse("mycmd -a=1 -b=k")
	cli.Option("-a").Int(0)
	got := cli.Option("-b").String("")
	if got != "k" {
		t.Fail()
	}
}

func Test_single_group_item_is_selected(t *testing.T) {
	cli := Parse("mycmd hello")
	phrases := cli.Group("Phrases", "PHRASE")
	phrases.New("hello", nil)
	phrases.Selected()
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_unknown_group_item(t *testing.T) {
	cli := Parse("mycmd car")
	nouns := cli.Group("Nouns", "NOUN")
	nouns.New("plane", nil)
	nouns.Selected()
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_required_argument(t *testing.T) {
	cli := Parse("mkdir")
	cli.Required("DIR")
	if cli.Ok() {
		t.Errorf("expected failure when required DIR is missing")
	}
}

func Test_optional_argument(t *testing.T) {
	cli := Parse("ls")
	cli.Optional("DIR")
	if !cli.Ok() {
		t.Error("unexpected:", cli.Error())
	}
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

func Test_usage_output_with_extended_docs(t *testing.T) {
	cli := NewParser()
	cli.args = []string{"adduser"}
	cli.Flag("-n, --dry-run")
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	cli.Option("-p, --password").String("")
	cli.Required("USERNAME").String("")

	var buf bytes.Buffer
	cli.WriteUsageTo(&buf)
	got := buf.String()
	if !strings.Contains(got, "one is generated") {
		t.Error("incomplete")
		t.Log(got)
	}
}

func Test_invalid_int_argument(t *testing.T) {
	cli := Parse("cmd -i=k")
	cli.Option("-i").Int(10)
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_undefined_option(t *testing.T) {
	cli := Parse("cmd -nosuch")
	if cli.Ok() {
		t.Error("should fail")
	}
}

func Test_stringer(t *testing.T) {
	exp := "mycmd -help -i=4"
	cli := Parse(exp)
	got := cli.String()
	if !strings.Contains(got, exp) {
		t.Error("\n", exp, "\n", got)
	}
}
