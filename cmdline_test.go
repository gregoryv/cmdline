package cmdline

import (
	"bytes"
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

func Test_error_handling_and_usage(t *testing.T) {
	cli := New("mycmd", "-i", "not-a-number")
	var exitCalled, usageCalled bool
	cli.exit = func(int) { exitCalled = true }
	cli.usage = func() { usageCalled = true }
	cli.Option("-i").Int(0)
	cli.CheckOptions()
	if !exitCalled {
		t.Error("Bad options should result in the exit func being called")
	}
	if !usageCalled {
		t.Error("Bad options should result in the usage func being called")
	}
}

func TestCommandLine_CheckOptions(t *testing.T) {
	var exitCalled, usageCalled bool
	cli := &CommandLine{
		args:  []string{"mycmd", "-h"},
		exit:  func(int) { exitCalled = true },
		usage: func() { usageCalled = true },
	}
	cli.CheckOptions()
	if !exitCalled {
		t.Error("-h flag should result in exit")
	}
	if !usageCalled {
		t.Error("-h flag should result in usage")
	}
}
