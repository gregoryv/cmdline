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
func (s *ShellOS) Getenv(key string) string { return os.Getenv(key) }

// Args returns os.Args
func (s *ShellOS) Args() []string { return os.Args }

// Getwd returns os.Getwd
func (s *ShellOS) Getwd() (string, error) { return os.Getwd() }

// Stdin returns os.Stdin
func (s *ShellOS) Stdin() io.Reader { return os.Stdin }

// Stdout returns os.Stdout
func (s *ShellOS) Stdout() io.Writer { return os.Stdout }

// Stderr returns os.Stderr
func (s *ShellOS) Stderr() io.Writer { return os.Stderr }

// Exit returns the given exit code
func (s *ShellOS) Exit(code int) { s.exit(code) }
func (s *ShellOS) Fatal(v ...interface{}) {
	log.Println(v...)
	s.exit(1)
}
