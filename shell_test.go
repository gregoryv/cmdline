package cmdline

import (
	"testing"
)

func Test_ShellOS(t *testing.T) {
	sh := NewShellOS()

	if sh.Getenv("PWD") == "" {
		t.Error(`sh.Getenv("PWD") failed `)
	}
	if len(sh.Args()) == 0 {
		t.Error("empty sh.Args")
	}
	wd, _ := sh.Getwd()
	if wd == "" {
		t.Error("empty sh.Getwd()")
	}
}

func Test_ShellOS_io(t *testing.T) {
	sh := NewShellOS()

	if sh.Stdin() == nil {
		t.Error("nil sh.Stdin")
	}
	if sh.Stdout() == nil {
		t.Error("nil sh.Stdout")
	}
	if sh.Stderr() == nil {
		t.Error("nil sh.Stderr")
	}
}

func Test_ShellOS_exits(t *testing.T) {
	sh := NewShellOS()

	// override the exit so we can test it
	sh.exit = func(v int) {
		if v != 1 {
			t.Error("got exit code", v)
		}
	}
	sh.Fatal()
	sh.Exit(1)
}
