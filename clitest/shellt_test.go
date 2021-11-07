package clitest

import (
	"strings"
	"testing"
)

func TestShellT(t *testing.T) {
	sh := NewShellT()
	defer sh.Cleanup()

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

func TestShellT_io(t *testing.T) {
	sh := NewShellT()
	defer sh.Cleanup()

	if sh.Stdin() == nil {
		t.Error("nil sh.Stdin")
	}
	if sh.Stdout() == nil {
		t.Error("nil sh.Stdout")
	}
	if sh.Stderr() == nil {
		t.Error("nil sh.Stderr")
	}
	sh.Fatal()
	sh.Exit(1)
}

func TestShellT_dumps(t *testing.T) {
	sh := NewShellT()
	defer sh.Cleanup()

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
