package cmdline

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
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
	cli.NamedArg("USERNAME").Strings()

	var buf bytes.Buffer
	cli.Usage().WriteTo(&buf)
	got := buf.String()
	if !strings.Contains(got, "one is generated") {
		t.Error("incomplete")
		t.Log(got)
	}
}

func ExampleParser_usageHiddenPassword() {
	cli := NewParser()
	cli.args = []string{"adduser"}
	_, opt := cli.Option("--uid").IntOpt(0)
	opt.Doc("If not given, one is generated")
	cli.Option("-p, --password",
		"minimum 8 chars",
		"hidden",
	).String("secret")

	cli.Usage().WriteTo(os.Stdout)
	// output:
	// Usage: adduser [OPTIONS]
	//
	// Options
	//     --uid : 0
	//         If not given, one is generated
	//
	//     -p, --password : "********"
	//         minimum 8 chars
}

func TestUsage_withoutGroups(t *testing.T) {
	cli := NewParser()
	cli.args = []string{"adduser"}
	_ = cli.Option("-u, --user-id").String("")
	_ = cli.Option("-p, --password", "hidden").String("")
	u := cli.Usage()
	u.Example("Add new user",
		"$ adduser -u john -p secret",
	)
	var buf bytes.Buffer
	cli.Usage().WriteTo(&buf)
	golden.Assert(t, buf.String())
}
