package cmdline

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func TestCommandLine_Usage(t *testing.T) {
	buf := bytes.NewBufferString("")
	cli := New("mycmd")
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc(
		"user id to set on the new account",
		"If not given, one is generated",
	)
	cli.Option("-u, --username").String("john")
	cli.Option("-p, --password").String("")
	cli.Flag("-n, --dry-run")
	cli.WriteUsageTo(buf)
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

func Test_error_handling(t *testing.T) {
	cases := []struct {
		cli     *CommandLine
		optFunc func(*CommandLine)
	}{
		{
			cli: newCli("mycmd -i 4 -unknown"),
			optFunc: func(cli *CommandLine) {
				cli.Option("-i").Int(0)
			},
		},
		{
			cli: newCli("mycmd -i notanumber"),
			optFunc: func(cli *CommandLine) {
				cli.Option("-i").Int(0)
			},
		},
		{
			cli:     newCli("mycmd -unknown"),
			optFunc: func(cli *CommandLine) {},
		},
	}
	for _, c := range cases {
		checkExit(t, c.cli, c.optFunc)
	}
}

func checkExit(t *testing.T, cli *CommandLine, fn func(*CommandLine)) {
	var exitCalled bool
	cli.exit = func(int) { exitCalled = true }
	out := bytes.NewBufferString("")
	cli.Output = out
	fn(cli)
	cli.CheckOptions()
	if !exitCalled {
		t.Error("Exit func not called for", cli)
		t.Log(cli.Args())
	}
}

func TestCommandLine_CheckOptions(t *testing.T) {
	var exitCalled, usageCalled bool
	cli := &CommandLine{
		args:   []string{"mycmd", "-h"},
		exit:   func(int) { exitCalled = true },
		usage:  func() { usageCalled = true },
		Output: ioutil.Discard,
	}
	cli.CheckOptions()
	if !exitCalled {
		t.Error("-h flag should result in exit")
	}
	if !usageCalled {
		t.Error("-h flag should result in usage")
	}
}

func TestCommandLine_Args(t *testing.T) {
	cli := &CommandLine{
		args: []string{"", "-h", "-i=3", "a", "b"},
	}
	cli.Option("-h").Bool()
	cli.Option("-i").Int(0)
	rest := cli.Args()
	if len(rest) != 2 {
		t.Errorf("Args did not return rest of arguments: %v", rest)
	}
}

func newCli(args string) *CommandLine {
	return New(strings.Split(args, " ")...)
}
