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
