package cmdline

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestNewShell(t *testing.T) {
	sh := NewShellOS()
	// override the exit so we can test it
	sh.exit = func(v int) {
		if v != 1 {
			t.Error("got exit code", v)
		}
	}
	assert := asserter.Wrap(t).Assert

	assert(sh.Getenv("PWD") != "").Error(`sh.Getenv("PWD") failed `)
	assert(len(sh.Args()) != 0).Error("empty sh.Args")

	wd, err := sh.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert(wd != "").Error("empty sh.Getwd()")
	assert(sh.Stdin() != nil).Error("nil sh.Stdin")
	assert(sh.Stdout() != nil).Error("nil sh.Stdout")
	assert(sh.Stderr() != nil).Error("nil sh.Stderr")
	sh.Fatal()
	sh.Exit(1)
}
