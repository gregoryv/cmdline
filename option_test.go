package cmdline

import (
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/cmdline/clitest"
)

func ExampleOption_Enum() {
	cli := NewParser()
	cli.Option("-a, --animal").Enum("snake", "snake", "bear", "goat")

	cli.Usage().WriteTo(os.Stdout)
	// output:
	//
	// Usage: mycmd [OPTIONS]
	//
	// Options
	//     -a, --animal : "snake" [snake bear goat]
}

func Test_ok_enum_single_value(t *testing.T) {
	cli := Parse(t, "cmd")
	cli.Option("-l, --letter").Enum("a")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_ok_enum(t *testing.T) {
	cli := Parse(t, "cmd -l a")
	cli.Option("-l, --letter").Enum("c", "a", "b")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_incorrect_enum(t *testing.T) {
	cli := Parse(t, "cmd -l x")
	cli.Option("-l, --letter").Enum("c", "a", "b", "c")
	if cli.Ok() {
		t.Error("x is not a valid enum")
	}
}

func Test_without_options(t *testing.T) {
	cli := Parse(t, "cmd")
	cli.Option("-b").String("default")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_non_options(t *testing.T) {
	cli := Parse(t, "cmd qwark")
	cli.Flag("-h")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_missing_string_option(t *testing.T) {
	cli := Parse(t, "cmd -a -b=1")
	got := cli.Option("-a").String("")
	cli.Option("-b").Int(0)
	if cli.Ok() {
		t.Error("should fail:", got)
	}
}

func Test_quoted_string_option(t *testing.T) {
	cli := Parse(t, `cmd -a "-b=1"`)
	got := cli.Option("-a").String("")
	if !cli.Ok() {
		t.Error(cli.Error(), got)
	}
}

func Test_default_uint8_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-a").Uint8(8)
	if got != 8 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uin16_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-a").Uint16(16)
	if got != 16 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uin32_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-a").Uint32(32)
	if got != 32 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uint_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-a").Uint(99)
	if got != 99 {
		t.Error("unexpected:", got)
	}
}

func Test_missing_uint_option(t *testing.T) {
	cli := Parse(t, "cmd -a")
	cli.Option("-a").Uint(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_bad_uint_option(t *testing.T) {
	cli := Parse(t, "cmd -a v")
	cli.Option("-a").Uint(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_default_int_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-a").Int(99)
	if got != 99 {
		t.Error("unexpected:", got)
	}
}

func Test_missing_int_option(t *testing.T) {
	cli := Parse(t, "cmd -a")
	cli.Option("-a").Int(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_default_bool_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-h").Bool()
	if got == true {
		t.Error("unexpected:", got)
	}
}

func Test_default_float64_option(t *testing.T) {
	cli := Parse(t, "cmd")
	got := cli.Option("-min").Float64(0.1)
	if got != 0.1 {
		t.Error("unexpected:", got)
	}
}

func Test_missing_float64_option(t *testing.T) {
	cli := Parse(t, "cmd -min")
	got := cli.Option("-min").Float64(0.1)
	if got != 0.1 {
		t.Error("unexpected:", got)
	}
}

func Test_bad_float64_option(t *testing.T) {
	cli := Parse(t, "cmd -min bad")
	got := cli.Option("-min").Float64(0.1)
	if got != 0.0 {
		t.Error("unexpected:", got)
	}
}

// Parse returns a parser from a string starting with the command
// followed by arguments.
func Parse(t *testing.T, str string) *Parser {
	p := NewParser()
	sh := clitest.NewShellT(strings.Split(str, " ")...)
	p.SetShell(sh)
	t.Cleanup(sh.Cleanup)
	return p
}
