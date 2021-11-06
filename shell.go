package cmdline

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

func (me *ShellOS) Getenv(key string) string { return os.Getenv(key) }
func (me *ShellOS) Args() []string           { return os.Args }
func (me *ShellOS) Getwd() (string, error)   { return os.Getwd() }
func (me *ShellOS) Stdin() io.Reader         { return os.Stdin }
func (me *ShellOS) Stdout() io.Writer        { return os.Stdout }
func (me *ShellOS) Stderr() io.Writer        { return os.Stderr }

// Exit returns the given exit code
func (me *ShellOS) Exit(code int) { me.exit(code) }
func (me *ShellOS) Fatal(v ...interface{}) {
	log.Println(v...)
	me.exit(1)
}

// NewShellT returns a shell with temporary working directory and
// buffered outputs, useful during testing.
// os.Chdir is called to change working directory to the temporary directory.
// The first argument should be name of command, just as in os.Args. If ommited
// /noname-tcmd is used. Temporary directory is based on that name.
func NewShellT(args ...string) *ShellT {
	if len(args) == 0 {
		args = []string{"/noname-tcmd"}
	}
	wd, err := ioutil.TempDir("", path.Base(args[0]))
	handleErr(err)
	origin, err := os.Getwd()
	handleErr(err)
	sh := &ShellT{
		Env: map[string]string{
			"PWD": wd,
		},
		args:     args,
		dir:      wd,
		origin:   origin,
		ExitCode: 0,
	}
	os.Chdir(sh.dir)
	return sh
}

var handleErr = func(err error) {
	if err != nil {
		panic(err)
	}
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

func (me *ShellT) Getenv(key string) (v string) {
	v, _ = me.Env[key]
	return
}

func (me *ShellT) Args() []string         { return me.args }
func (me *ShellT) Getwd() (string, error) { return me.dir, nil }
func (me *ShellT) Stdin() io.Reader       { return &me.In }
func (me *ShellT) Stdout() io.Writer      { return &me.Out }
func (me *ShellT) Stderr() io.Writer      { return &me.Err }

// Exit sets ExitCode
func (me *ShellT) Exit(code int) {
	me.ExitCode = code
	os.Chdir(me.origin)
}

// Fatal logs the given values and calls the Exit method
func (me *ShellT) Fatal(v ...interface{}) {
	log.Println(v...)
	os.Chdir(me.origin)
	me.Exit(1)
}

// Cleanup removes temporary directory and restores the working
// directory.
func (me *ShellT) Cleanup() {
	os.Chdir(me.origin)
	os.RemoveAll(me.dir)
}

// Dump returns a dump of the command, see DumpTo
func (me *ShellT) Dump() string {
	var b strings.Builder
	me.DumpTo(&b)
	return b.String()
}

// DumpTo writes argument, stdout and stderr if any to the given writer
func (me *ShellT) DumpTo(w io.Writer) error {
	p, err := nexus.NewPrinter(w)
	p.Print("$ ")
	p.Print(strings.Join(me.Args(), " "))
	p.Println()
	io.Copy(p, &me.Out)
	p.Println()
	p.Print("exit ", me.ExitCode)

	if me.Err.Len() > 0 {
		p.Println()
		p.Println("STDERR:")
		io.Copy(p, &me.Err)
	}
	return *err
}
