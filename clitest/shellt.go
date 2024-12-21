package clitest

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gregoryv/nexus"
)

// NewShellT returns a shell with temporary working directory and
// buffered outputs, useful during testing.
// os.Chdir is called to change working directory to the temporary directory.
// The first argument should be name of command, just as in os.Args. If ommited
// /noname-tcmd is used. Temporary directory is based on that name.
func NewShellT(args ...string) *ShellT {
	if len(args) == 0 {
		args = []string{"/noname-tcmd"}
	}
	origin, _ := os.Getwd() // store for later
	wd, _ := ioutil.TempDir("", path.Base(args[0]))
	os.Chdir(wd)
	sh := &ShellT{
		Env: map[string]string{
			"PWD": wd,
		},
		args:     args,
		dir:      wd,
		origin:   origin,
		ExitCode: 0,
	}
	return sh
}

type ShellT struct {
	Env      map[string]string
	Out      bytes.Buffer // Stdout
	Err      bytes.Buffer // Stderr
	In       bytes.Buffer // Stdin
	ExitCode int          // Set by method Exit

	args   []string
	dir    string
	origin string
}

func (s *ShellT) Getenv(key string) (v string) {
	v, _ = s.Env[key]
	return
}

func (s *ShellT) Args() []string         { return s.args }
func (s *ShellT) Getwd() (string, error) { return s.dir, nil }
func (s *ShellT) Stdin() io.Reader       { return &s.In }
func (s *ShellT) Stdout() io.Writer      { return &s.Out }
func (s *ShellT) Stderr() io.Writer      { return &s.Err }

// Exit sets ExitCode
func (s *ShellT) Exit(code int) {
	s.ExitCode = code
	os.Chdir(s.origin)
}

// Fatal logs the given values and calls the Exit method
func (s *ShellT) Fatal(v ...interface{}) {
	log.Println(v...)
	os.Chdir(s.origin)
	s.Exit(1)
}

// Cleanup removes temporary directory and restores the working
// directory.
func (s *ShellT) Cleanup() {
	os.Chdir(s.origin)
	os.RemoveAll(s.dir)
}

// Dump returns a dump of the command, see DumpTo
func (s *ShellT) Dump() string {
	var b strings.Builder
	s.DumpTo(&b)
	return b.String()
}

// DumpTo writes argument, stdout and stderr if any to the given writer
func (s *ShellT) DumpTo(w io.Writer) error {
	p, err := nexus.NewPrinter(w)
	p.Print("$ ")
	p.Print(strings.Join(s.Args(), " "))
	p.Println()
	io.Copy(p, &s.Out)
	p.Println()
	p.Print("exit ", s.ExitCode)

	if s.Err.Len() > 0 {
		p.Println()
		p.Println("STDERR:")
		io.Copy(p, &s.Err)
	}
	return *err
}
