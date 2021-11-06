package cmdline

import (
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestNewShell(t *testing.T) {
	sh, _ := NewShell().(*ossh)
	sh.exit = func(v int) {
		if v != 1 {
			t.Error("got exit code", v)
		}
	}
	testShell(t, sh)
}

func TestNewTestShell(t *testing.T) {
	sh := NewTestShell()
	defer sh.Cleanup()

	testShell(t, sh)

	// dump contains exit code
	got := sh.Dump()
	if !strings.Contains(got, "exit ") {
		t.Error(got)
	}

	// dump can contain content of stderr
	sh.Err.WriteString("hello")
	got = sh.Dump()
	if !strings.Contains(got, "hello") {
		t.Error(got)
	}

}

func testShell(t *testing.T, cmd Shell) {
	t.Helper()
	assert := asserter.Wrap(t).Assert

	assert(cmd.Getenv("PWD") != "").Error(`cmd.Getenv("PWD") failed `)
	assert(len(cmd.Args()) != 0).Error("empty cmd.Args")

	wd, _ := cmd.Getwd()
	assert(wd != "").Error("empty cmd.Getwd()")
	assert(cmd.Stdin() != nil).Error("nil cmd.Stdin")
	assert(cmd.Stdout() != nil).Error("nil cmd.Stdout")
	assert(cmd.Stderr() != nil).Error("nil cmd.Stderr")
	cmd.Fatal()
	cmd.Exit(1)
}
