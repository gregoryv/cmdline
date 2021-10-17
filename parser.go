package cmdline

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// NewBasicParser returns a parser including help options -h, --help.
func NewBasicParser() *Basic {
	return &Basic{Parser: NewParser()}
}

type Basic struct {
	*Parser

	sync.Once // used to parse the help flag only once
	help      bool
}

// Parse checks for errors or if the help flag is given writes usage
// to os.Stdout
func (me *Basic) Parse() {
	if !me.Ok() {
		log.Println(me.Error())
		log.Println("Try -h or --help, for more information")
		me.exit(1)
	}
	me.Once.Do(me.helpFlag)
	if me.help {
		me.WriteUsageTo(os.Stdout)
		me.exit(0)
	}
}

func (me *Basic) WriteUsageTo(w io.Writer) {
	me.Once.Do(me.helpFlag)
	me.Parser.WriteUsageTo(w)
}

func (me *Basic) helpFlag() {
	me.help = me.Flag("-h, --help")
}

// ----------------------------------------

// NewParser returns a parser. First argument must be the command
// name.
func NewParser() *Parser {
	return &Parser{
		args:    os.Args,
		options: make([]*Option, 0),
		groups:  make([]*Group, 0),
		envMap:  os.Getenv,
		exit:    os.Exit,
	}
}

// Parser groups arguments for option parsing and usage.
type Parser struct {
	args      []string // including command name as first element
	options   []*Option
	arguments []*Argument // required

	groups []*Group

	envMap func(string) string
	exit   func(int)
}

// Parse checks parsing errors and exits on errors
func (me *Parser) Parse() {
	if !me.Ok() {
		log.Println(me.Error())
		me.exit(1)
	}
}

func (me *Parser) SetArgs(args ...string) { me.args = args }
func (me *Parser) SetExit(v func(int))    { me.exit = v }

func (me *Parser) Group(title, name string, items ...*Item) *Group {
	return me.group(title, name, me.Optional(name).String(""), items)
}

func (me *Parser) group(title, name, v string, items []*Item) *Group {
	grp := &Group{
		name:  name,
		args:  me.Args(),
		title: title,
		v:     v,
		items: items,
	}
	err := me.addGroup(grp)
	if err != nil {
		panic(fmt.Sprintf("duplicate group %q", title))
	}

	return grp
}

func (me *Parser) addGroup(grp *Group) error {
	for _, existing := range me.groups {
		if existing.Title() == grp.Title() {
			return fmt.Errorf("group %q already exists", grp.Title())
		}
	}
	me.groups = append(me.groups, grp)
	return nil
}

type Group struct {
	name string
	args []string // needed for parsing extra options

	title string
	v     string
	items []*Item

	err error
}

func (me *Group) New(name string, any interface{}) *Item {
	item := &Item{
		Name:   name,
		Loader: any,
	}
	me.items = append(me.items, item)
	return item
}

// Selected returns the matching item. Defaults to the first in the group.
func (me *Group) Selected() interface{} {
	i := me.items[0]
	if me.v != "" {
		var found bool
		i, found = me.find(me.v)
		if !found {
			me.err = fmt.Errorf("invalid %s", me.name)
			return nil
		}
	}
	extra := NewParser()
	extra.args = append([]string{me.title}, me.args...)
	sel := i.Load(extra)
	me.err = extra.Error()
	return sel
}

func (me *Group) Title() string  { return me.title }
func (me *Group) Items() []*Item { return me.items }

// Find returns the named Item or nil if not found.
func (me *Group) find(name string) (*Item, bool) {
	for _, item := range me.items {
		if item.Name == name {
			return item, true
		}
	}
	return nil, false
}

// ----------------------------------------

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
	if len(me.groups) == 0 { // as groups are selected with non option argument
		for _, arg := range me.Args() {
			if isOption(arg) {
				return fmt.Errorf("Unknown option: %v", arg)
			}
		}
	}
	return nil
}

func (me *Parser) parseFailed() error {
	var err error
	setErr := func(e error) {
		if err != nil {
			return
		}
		err = e
	}
	for _, opt := range me.options {
		setErr(opt.err)
	}
	for _, arg := range me.arguments {
		setErr(arg.err)
	}
	for _, grp := range me.groups {
		setErr(grp.err)
	}
	return err
}

// Option returns a new option with the given names.
// Names should be a comma separated string, e.g.
//   -n, --dry-run
//
// You can also include an environment variable to use as default
// value in the names, e.g.
//   -t, --token, $COGNITO_TOKEN
func (me *Parser) Option(names string, doclines ...string) *Option {
	opt := NewOption(names, me.args[1:]...)
	opt.envMap = me.envMap
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

	indent := "    "
	for _, grp := range me.groups {
		fmt.Fprintln(w, grp.Title())
		first := grp.Items()[0]
		writeItem(w, first, me.args, indent, true)
		for _, item := range grp.Items()[1:] {
			writeItem(w, item, me.args, indent, false)
		}
	}
}

func writeItem(w io.Writer, me *Item, args []string, indent string, dflt bool) {
	if dflt {
		fmt.Fprintf(w, "%s%s (default)\n", indent, me.Name)
	} else {
		fmt.Fprintf(w, "%s%s\n", indent, me.Name)
	}
	extra := NewParser()
	extra.args = args
	me.Load(extra)
	extra.writeOptionsTo(w, indent)
}

// WriteOptionsTo writes the Options section to the given writer.
func (me *Parser) WriteOptionsTo(w io.Writer) {
	me.writeOptionsTo(w, "")
}

func (me *Parser) writeOptionsTo(w io.Writer, indent string) {
	for _, opt := range me.options {
		writeOptionTo(w, opt, indent)
	}
	if len(me.options) > 0 {
		fmt.Fprintln(w)
	}
}

func writeOptionTo(w io.Writer, opt *Option, indent string) {
	var def string
	if opt.defaultValue != "" {
		def = fmt.Sprintf(" : %v", opt.defaultValue)
	}
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

// Optional returns an optional  named argument.
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
func (me *Argument) String(def string) string {
	if me.v == "" {
		return def
	}
	return me.v
}
