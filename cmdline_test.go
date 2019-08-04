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
	var called bool
	cli.exit = func(int) { called = true }
	cli.Option("-i").Int(0)
	cli.CheckOptions()
	if !called {
		t.Error("Bad options should result in the exit func being called")
	}
}

func TestCommandLine_CheckOptions(t *testing.T) {
	cli := New("mycmd", "-h")
	var called bool
	cli.exit = func(int) { called = true }
	cli.CheckOptions()
	if !called {
		t.Error("-h flag should result in exit")
	}
}
