package cmdline

import (
	"bytes"
	"strings"
	"testing"
)

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
	cli.Usage().WriteTo(&buf)
	got := buf.String()
	if !strings.Contains(got, "one is generated") {
		t.Error("incomplete")
		t.Log(got)
	}
}
