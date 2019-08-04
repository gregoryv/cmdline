package cmdline

import (
	"fmt"
	"io"
	"strings"
)

type CommandLine struct {
	args    []string // including command name as first element
	usage   func()
	exit    func(int)
	options []*Option
}

func New(args ...string) *CommandLine {
	if len(args) == 0 {
		panic("New() missing args")
	}
	return &CommandLine{
		args:    args,
		usage:   noUsage,
		exit:    noExit,
		options: make([]*Option, 0),
	}
}

func noUsage()   {}
func noExit(int) {}

func (cli *CommandLine) string(names string) (string, bool) {
	for i, arg := range cli.args {
		if namesMatch(names, arg) {
			isLast := len(cli.args)-1 == i
			if isLast {
				return "", true
			}
			return cli.args[i+1], true
		}
	}
	return "", false
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

func (cli *CommandLine) Option(names string) *Option {
	opt := &Option{
		args:  cli.args[1:],
		names: names,
	}
	cli.options = append(cli.options, opt)
	return opt
}

func (cli *CommandLine) Flag(name string) bool {
	val, _ := cli.Option(name).BoolOpt()
	return val
}

func (cli *CommandLine) WriteUsageTo(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [OPTIONS]\n\n", cli.args[0])
	cli.WriteOptionsTo(w)
}

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
