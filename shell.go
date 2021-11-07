package cmdline

import (
	"io"
	"log"
	"os"
)

// Shell defines a command line execution context.
type Shell interface {
	Getenv(string) string
	Args() []string
	Getwd() (string, error)
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Exit(code int)
	Fatal(v ...interface{})
}

// NewShellOS returns a shell from the os values.
func NewShellOS() *ShellOS {
	return &ShellOS{
		exit: os.Exit,
	}
}

// ShellOS uses package os
type ShellOS struct {
	exit func(code int)
}

// Getenv returns os.Getenv
func (me *ShellOS) Getenv(key string) string { return os.Getenv(key) }

// Args returns os.Args
func (me *ShellOS) Args() []string { return os.Args }

// Getwd returns os.Getwd
func (me *ShellOS) Getwd() (string, error) { return os.Getwd() }

// Stdin returns os.Stdin
func (me *ShellOS) Stdin() io.Reader { return os.Stdin }

// Stdout returns os.Stdout
func (me *ShellOS) Stdout() io.Writer { return os.Stdout }

// Stderr returns os.Stderr
func (me *ShellOS) Stderr() io.Writer { return os.Stderr }

// Exit returns the given exit code
func (me *ShellOS) Exit(code int) { me.exit(code) }
func (me *ShellOS) Fatal(v ...interface{}) {
	log.Println(v...)
	me.exit(1)
}
