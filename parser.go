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
	return &Parser{
		args:    args,
		options: make([]*Option, 0),
		groups:  make([]*Group, 0),
	}
}

// Parser groups arguments for option parsing and usage.
type Parser struct {
	args      []string // including command name as first element
	options   []*Option
	arguments []*Argument // required

	groups []*Group
}

func (me *Parser) Group(name string, v ...Item) (*Group, error) {
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
func (me *Parser) Ok() bool {
	return me.Error() == nil
}

// Error returns first error of the given options.
func (me *Parser) Error() error {
	err := me.parseFailed()
	if err != nil {
		return err
	}
	for _, arg := range me.Args() {
		if isOption(arg) {
			return fmt.Errorf("Unknown option: %v", arg)
		}
	}
	return nil
}

func (me *Parser) parseFailed() error {
	for _, opt := range me.options {
		if opt.err != nil {
			return opt.err
		}
	}
	for _, arg := range me.arguments {
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
func (me *Parser) Option(names string, doclines ...string) *Option {
	opt := NewOption(names, me.args[1:]...)
	opt.doc = doclines
	me.options = append(me.options, opt)
	return opt
}

// Flag is short for Option(name).Bool()
func (me *Parser) Flag(name string) bool {
	val, _ := me.Option(name).BoolOpt()
	return val
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//   Usage: COMMAND [OPTIONS] ARGUMENTS...
func (me *Parser) WriteUsageTo(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [OPTIONS]", me.args[0])
	for _, arg := range me.arguments {
		if arg.required {
			fmt.Fprintf(w, " %s", arg.name)
			continue
		}
		fmt.Fprintf(w, " [%s]", arg.name)
	}
	fmt.Fprint(w, "\n\n")
	fmt.Fprintln(w, "Options")
	me.WriteOptionsTo(w)

	for _, grp := range me.groups {
		indent := "    "
		for _, i := range grp.Items() {
			fmt.Fprintln(w)
			fmt.Fprintln(w, grp.Name())
			fmt.Fprintf(w, "%s%s\n", indent, i.Name())
			extra := NewParser(os.Args...)
			i.ExtraOptions(extra)
			extra.writeOptionsTo(w, indent)
		}
	}
}

// WriteOptionsTo writes the Options section to the given writer.
func (me *Parser) WriteOptionsTo(w io.Writer) {
	me.writeOptionsTo(w, "")
}

func (me *Parser) writeOptionsTo(w io.Writer, indent string) {
	for _, opt := range me.options {
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
func (me *Parser) Args() []string {
	rest := make([]string, 0)
	for i, arg := range me.args[1:] {
		//		fmt.Println("a:", arg)
		if !me.wasMatched(i) {
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

func (me *Parser) wasMatched(i int) bool {
	for _, opt := range me.options {
		if opt.argIndex == i || opt.valIndex == i {
			return true
		}
	}
	return false
}

func (me *Parser) String() string {
	return fmt.Sprintf("Parser: %s", strings.Join(me.args, " "))
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

// String returns the value of this argument
func (me *Argument) String() string { return me.v }
