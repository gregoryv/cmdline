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

// NewShell returns a shell from the os values.
func NewShell() Shell {
	return &ossh{
		exit: os.Exit,
	}
}

type ossh struct {
	exit func(code int)
}

func (me *ossh) Getenv(key string) string { return os.Getenv(key) }
func (me *ossh) Args() []string           { return os.Args }
func (me *ossh) Getwd() (string, error)   { return os.Getwd() }
func (me *ossh) Stdin() io.Reader         { return os.Stdin }
func (me *ossh) Stdout() io.Writer        { return os.Stdout }
func (me *ossh) Stderr() io.Writer        { return os.Stderr }

// Exit returns the given exit code
func (me *ossh) Exit(code int) { me.exit(code) }
func (me *ossh) Fatal(v ...interface{}) {
	log.Println(v...)
	me.exit(1)
}

// NewTestShell returns a shell with temporary working directory and
// buffered outputs, useful during testing.
// os.Chdir is called to change working directory to the temporary directory.
// The first argument should be name of command, just as in os.Args. If ommited
// /noname-tcmd is used. Temporary directory is based on that name.
func NewTestShell(args ...string) *TestShell {
	if len(args) == 0 {
		args = []string{"/noname-tcmd"}
	}
	wd, err := ioutil.TempDir("", path.Base(args[0]))
	handleErr(err)
	origin, err := os.Getwd()
	handleErr(err)
	sh := &TestShell{
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

type TestShell struct {
	Env      map[string]string
	Out      bytes.Buffer // Stdout
	Err      bytes.Buffer // Stderr
	In       bytes.Buffer // Stdin
	ExitCode int          // Set by method Exit

	args   []string
	dir    string
	origin string
}

func (me *TestShell) Getenv(key string) (v string) {
	v, _ = me.Env[key]
	return
}

func (me *TestShell) Args() []string         { return me.args }
func (me *TestShell) Getwd() (string, error) { return me.dir, nil }
func (me *TestShell) Stdin() io.Reader       { return &me.In }
func (me *TestShell) Stdout() io.Writer      { return &me.Out }
func (me *TestShell) Stderr() io.Writer      { return &me.Err }

// Exit sets ExitCode
func (me *TestShell) Exit(code int) { me.ExitCode = code }

// Fatal logs the given values and calls the Exit method
func (me *TestShell) Fatal(v ...interface{}) {
	log.Println(v...)
	me.Exit(1)
}

// Cleanup removes temporary directory and restores the working
// directory.
func (me *TestShell) Cleanup() {
	os.Chdir(me.origin)
	os.RemoveAll(me.dir)
}

// Dump returns a dump of the command, see DumpTo
func (me *TestShell) Dump() string {
	var b strings.Builder
	me.DumpTo(&b)
	return b.String()
}

// DumpTo writes argument, stdout and stderr if any to the given writer
func (me *TestShell) DumpTo(w io.Writer) error {
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
