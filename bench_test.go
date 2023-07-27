package cmdline

import (
	"flag"
	"strings"
	"testing"

	"github.com/gregoryv/cmdline/clitest"
)

var txt = `adduser -vD -u john --password=secret -uid 101 ` +
	`".. description .." "More text." word -- -v -x`

var sameArgs = strings.Split(txt, " ")

func BenchmarkParse(b *testing.B) {
	sh := clitest.NewShellT(sameArgs...)
	for i := 0; i < b.N; i++ {
		cli := NewBasicParser()
		cli.SetShell(sh)
		_ = cli.Flag("-v")
		_ = cli.Flag("-d")
		cli.Parse()
	}
}

func BenchmarkFlagSet_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		v := fs.Bool("v", false, "")
		d := fs.Bool("d", false, "")
		fs.Parse(sameArgs)

		_ = v
		_ = d
	}
}
