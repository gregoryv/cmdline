package cmdline

import (
	"testing"
)

func Test_without_options(t *testing.T) {
	cli := Parse("cmd")
	cli.Option("-b").String("default")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_non_options(t *testing.T) {
	cli := Parse("cmd qwark")
	cli.Flag("-h")
	if !cli.Ok() {
		t.Error(cli.Error())
	}
}

func Test_missing_string_option(t *testing.T) {
	cli := Parse("cmd -a -b=1")
	got := cli.Option("-a").String("")
	cli.Option("-b").Int(0)
	if cli.Ok() {
		t.Error("should fail:", got)
	}
}

func Test_quoted_string_option(t *testing.T) {
	cli := Parse(`cmd -a "-b=1"`)
	got := cli.Option("-a").String("")
	if !cli.Ok() {
		t.Error(cli.Error(), got)
	}
}

func Test_default_uint8_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-a").Uint8(8)
	if got != 8 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uin16_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-a").Uint16(16)
	if got != 16 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uin32_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-a").Uint32(32)
	if got != 32 {
		t.Error("unexpected:", got)
	}
}

func Test_default_uint_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-a").Uint(99)
	if got != 99 {
		t.Error("unexpected:", got)
	}
}

func Test_missing_uint_option(t *testing.T) {
	cli := Parse("cmd -a")
	cli.Option("-a").Uint(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_bad_uint_option(t *testing.T) {
	cli := Parse("cmd -a v")
	cli.Option("-a").Uint(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_default_int_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-a").Int(99)
	if got != 99 {
		t.Error("unexpected:", got)
	}
}

func Test_missing_int_option(t *testing.T) {
	cli := Parse("cmd -a")
	cli.Option("-a").Int(0)
	if cli.Ok() {
		t.Fail()
	}
}

func Test_default_bool_option(t *testing.T) {
	cli := Parse("cmd")
	got := cli.Option("-h").Bool()
	if got == true {
		t.Error("unexpected:", got)
	}
}
