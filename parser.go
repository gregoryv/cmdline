package cmdline

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Parse returns a parser from a string starting with the command
// followed by arguments.
func Parse(str string) *Parser {
	return NewParser(strings.Split(str, " ")...)
}

// NewParser returns a parser, usually called with
// cmdline.New(os.Args...).  First argument must be the command name
func NewParser(args ...string) *Parser {
	if len(args) == 0 {
		panic("New() missing args")
	}
	cli := &Parser{
		args:    args,
		options: make([]*Option, 0),
		groups:  make([]*Group, 0),
	}
	return cli
}

// Parser groups arguments for option parsing and usage.
type Parser struct {
	args      []string // including command name as first element
	options   []*Option
	arguments []*Argument // required

	groups []*Group
}

func (me *Parser) Group(name string, v ...Action) (*Group, error) {
	grp := NewGroup(name, v...)
	return grp, me.AddGroup(grp)
}

func (me *Parser) AddGroup(grp *Group) error {
	for _, existing := range me.groups {
		if existing.Name() == grp.Name() {
			return fmt.Errorf("group %q already exists", grp.Name())
		}
	}
	me.groups = append(me.groups, grp)
	return nil
}

// Ok returns true if no parsing error occured
func (cli *Parser) Ok() bool {
	return cli.Error() == nil
}

// Error returns first error of the given options.
func (cli *Parser) Error() error {
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

func (cli *Parser) parseFailed() error {
	for _, opt := range cli.options {
		if opt.err != nil {
			return opt.err
		}
	}
	for _, arg := range cli.arguments {
		if arg.err != nil {
			return arg.err
		}
	}
	return nil
}

// Option returns a new option with the given names.
// Names should be a comma separated string, e.g.
//   -n, --dry-run
//
func (cli *Parser) Option(names string, doclines ...string) *Option {
	opt := NewOption(names, cli.args[1:]...)
	opt.doc = doclines
	cli.options = append(cli.options, opt)
	return opt
}

// Flag is short for Option(name).Bool()
func (cli *Parser) Flag(name string) bool {
	val, _ := cli.Option(name).BoolOpt()
	return val
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//   Usage: COMMAND [OPTIONS] ARGUMENTS...
func (cli *Parser) WriteUsageTo(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [OPTIONS]", cli.args[0])
	for _, arg := range cli.arguments {
		if arg.required {
			fmt.Fprintf(w, " %s", arg.name)
			continue
		}
		fmt.Fprintf(w, " [%s]", arg.name)
	}
	fmt.Fprint(w, "\n\n")
	fmt.Fprintln(w, "Options")
	cli.WriteOptionsTo(w)

	for _, grp := range cli.groups {
		indent := "    "
		for _, action := range grp.actions {
			fmt.Fprintln(w)
			fmt.Fprintln(w, grp.Name())
			fmt.Fprintf(w, "%s%s\n", indent, action.Name())
			extra := NewParser(os.Args...)
			action.ExtraOptions(extra)
			extra.writeOptionsTo(w, indent)
		}
	}
}

// WriteOptionsTo writes the Options section to the given writer.
func (cli *Parser) WriteOptionsTo(w io.Writer) {
	cli.writeOptionsTo(w, "")
}

func (cli *Parser) writeOptionsTo(w io.Writer, indent string) {
	for _, opt := range cli.options {
		def := fmt.Sprintf(" : %v", opt.defaultValue)
		if opt.quoteValue {
			def = fmt.Sprintf(" : %q", opt.defaultValue)
		}
		fmt.Fprintf(w, "%s    %s%s\n", indent, opt.names, def)
		if len(opt.doc) > 0 {
			for _, line := range opt.doc {
				fmt.Fprintf(w, "%s        %s\n", indent, line)
			}
			fmt.Fprintln(w)
		}
	}
}

// Args returns arguments not matched by any of the options
func (cli *Parser) Args() []string {
	rest := make([]string, 0)
	for i, arg := range cli.args[1:] {
		//		fmt.Println("a:", arg)
		if !cli.wasMatched(i) {
			rest = append(rest, arg)
		}
	}
	return rest
}

// Argn returns the n:th argument of remaining arguments starting at 0.
func (me *Parser) Argn(n int) string {
	rest := me.Args()
	if n < len(rest) {
		return rest[n]
	}
	return ""
}

func (cli *Parser) wasMatched(i int) bool {
	for _, opt := range cli.options {
		if opt.argIndex == i || opt.valIndex == i {
			return true
		}
	}
	return false
}

func (cli *Parser) String() string {
	return fmt.Sprintf("Parser: %s", strings.Join(cli.args, " "))
}

// Required returns a required named argument.
func (me *Parser) Required(name string) *Argument {
	arg := &Argument{
		name:     name,
		v:        me.Argn(len(me.arguments)),
		required: true,
	}
	if arg.v == "" {
		arg.err = fmt.Errorf("missing %s", name)
	}
	me.arguments = append(me.arguments, arg)
	return arg
}

// Optional returns a required named argument.
func (me *Parser) Optional(name string) *Argument {
	arg := &Argument{
		name: name,
		v:    me.Argn(len(me.arguments)),
	}
	me.arguments = append(me.arguments, arg)
	return arg
}

type Argument struct {
	name     string
	v        string
	err      error
	required bool
}

// String
func (me *Argument) String() string {
	return me.v
}
