package cmdline

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// CommandLine groups arguments for option parsing and usage.
type CommandLine struct {
	args    []string // including command name as first element
	usage   func()
	exit    func(int)
	options []*Option
}

// New returns a CommandLine, usually called with cmdline.New(os.Args...).
// First argument must be the command name
func New(args ...string) *CommandLine {
	if len(args) == 0 {
		panic("New() missing args")
	}

	cli := &CommandLine{
		args:    args,
		exit:    os.Exit,
		options: make([]*Option, 0),
	}
	cli.usage = func() {
		cli.WriteUsageTo(os.Stderr)
	}
	return cli
}

func namesMatch(names, arg string) bool {
	if !isOption(arg) {
		return false
	}
	return strings.Index(names, arg) >= 0
}

func isOption(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

// CheckOptions exits if any of the given options are incorrect.
func (cli *CommandLine) CheckOptions() {
	err := cli.parseFailed()
	if err != nil {
		cli.usage()
		cli.exit(1)
	}
}

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
	opt := &Option{
		args:  cli.args[1:],
		names: names,
	}
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
	fmt.Fprintf(w, "Usage: %s [OPTIONS]\n\n", cli.args[0])
	cli.WriteOptionsTo(w)
}

// WriteOptionsTo writes the Options section to the given writer.
func (cli *CommandLine) WriteOptionsTo(w io.Writer) {
	fmt.Fprintln(w, "Options")
	for _, opt := range cli.options {
		indent := "\t\t"
		if len(opt.names) < 10 {
			indent = "\t\t\t"
		}
		def := fmt.Sprintf("%s(default: %s)", indent, opt.defaultValue)
		fmt.Fprintf(w, "    %s%s\n", opt.names, def)
		if len(opt.doc) > 0 {
			for _, line := range opt.doc {
				fmt.Fprintln(w, "\t", line)
			}
			fmt.Fprintln(w)
		}
	}
}
