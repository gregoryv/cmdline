package cmdline

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Parse returns a command line from a string starting with the
// command followed by arguments.
func Parse(str string) *CommandLine {
	return New(strings.Split(str, " ")...)
}

// New returns a CommandLine, usually called with
// cmdline.New(os.Args...).  First argument must be the command name
func New(args ...string) *CommandLine {
	if len(args) == 0 {
		panic("New() missing args")
	}
	cli := &CommandLine{
		args:    args,
		options: make([]*Option, 0),
	}
	return cli
}

// CommandLine groups arguments for option parsing and usage.
type CommandLine struct {
	args      []string // including command name as first element
	options   []*Option
	arguments []*Argument // required

	Help bool // set to true if -h or --help is given
}

// Error returns first error of the given options.
func (cli *CommandLine) Error() error {
	err := cli.parseFailed()
	if err != nil {
		return err
	}
	for _, arg := range cli.Args() {
		if isOption(arg) {
			return fmt.Errorf("Unknown option: %v", arg)
		}
	}
	return nil
}

var ErrHelp = errors.New("Help requested")

func (cli *CommandLine) parseFailed() error {
	for _, opt := range cli.options {
		if opt.err != nil {
			return opt.err
		}
	}
	return nil
}

// Option returns a new option with the given names.
// Names should be a comma separated string, e.g.
//   -n, --dry-run
//
func (cli *CommandLine) Option(names string) *Option {
	opt := NewOption(names, cli.args[1:]...)
	cli.options = append(cli.options, opt)
	return opt
}

// Flag is short for Option(name).Bool()
func (cli *CommandLine) Flag(name string) bool {
	val, _ := cli.Option(name).BoolOpt()
	return val
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//   Usage: COMMAND [OPTIONS]
func (cli *CommandLine) WriteUsageTo(w io.Writer) {
	required := make([]string, len(cli.arguments))
	for i, argument := range cli.arguments {
		required[i] = argument.name
	}
	var req string
	if len(required) > 0 {
		req = fmt.Sprintf(" %s", strings.Join(required, " "))
	}
	fmt.Fprintf(w, "Usage: %s [OPTIONS]%s\n\n", cli.args[0], req)
	cli.WriteOptionsTo(w)
}

// WriteOptionsTo writes the Options section to the given writer.
func (cli *CommandLine) WriteOptionsTo(w io.Writer) {
	fmt.Fprintln(w, "Options")
	for _, opt := range cli.options {
		def := fmt.Sprintf(" : %v", opt.defaultValue)
		if opt.quoteValue {
			def = fmt.Sprintf(" : %q", opt.defaultValue)
		}
		fmt.Fprintf(w, "    %s%s\n", opt.names, def)
		if len(opt.doc) > 0 {
			for _, line := range opt.doc {
				fmt.Fprintln(w, "\t", line)
			}
			fmt.Fprintln(w)
		}
	}
}

// Args returns arguments not matched by any of the options
func (cli *CommandLine) Args() []string {
	rest := make([]string, 0)
	for i, arg := range cli.args[1:] {
		//		fmt.Println("a:", arg)
		if !cli.wasMatched(i) {
			rest = append(rest, arg)
		}
	}
	return rest
}

// Argn returns the n:th argument of remaining arguments
func (me *CommandLine) Argn(n int) string {
	rest := me.Args()
	if n < len(rest) {
		return rest[n]
	}
	return ""
}

func (cli *CommandLine) wasMatched(i int) bool {
	for _, opt := range cli.options {
		if opt.argIndex == i || opt.valIndex == i {
			return true
		}
	}
	return false
}

func (cli *CommandLine) String() string {
	return fmt.Sprintf("CommandLine: %s", strings.Join(cli.args, " "))
}

// NeedArg returns a named argument.
func (me *CommandLine) NeedArg(name string, n int) *Argument {
	arg := &Argument{
		name: name,
		v:    me.Argn(n),
	}
	me.arguments = append(me.arguments, arg)
	return arg
}

type Argument struct {
	name string
	v    string
}

// String
func (me *Argument) String() string {
	return me.v
}
