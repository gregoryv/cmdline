// Package cmdline provides a way to parse command line arguments.
package cmdline

import (
	"fmt"
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

	defineHelp sync.Once // used to parse the help flag only once
	help       bool
}

// Parse checks for errors or if the help flag is given writes usage
// to os.Stdout
func (me *Basic) Parse() {
	me.defineHelp.Do(me.helpFlag)

	switch {
	case me.help:
		me.Usage().WriteTo(os.Stdout)
		me.sh.Exit(0)

	case !me.Ok():
		fmt.Println(me.Error())
		fmt.Println("Try -h or --help, for more information")
		me.sh.Exit(1)
	}
}

// Usage returns the usage for further documentation. If Parse method
// has not been called, it adds the help flag.
func (me *Basic) Usage() *Usage {
	me.defineHelp.Do(me.helpFlag)
	return me.Parser.Usage()
}

func (me *Basic) helpFlag() {
	me.help = me.Parser.Flag("-h, --help")
}

// ----------------------------------------
var DefaultShell = NewShellOS()

// NewParser returns a parser using the DefaultShell
func NewParser() *Parser {
	sh := DefaultShell
	p := &Parser{
		sh:      sh,
		args:    sh.Args(),
		options: make([]*Option, 0),
		groups:  make([]*Group, 0),
		envMap:  sh.Getenv,
	}
	p.usage = &Usage{Parser: p}
	return p
}

// Parser groups arguments for option parsing and usage.
type Parser struct {
	sh Shell

	args      []string // including command name as first element
	options   []*Option
	arguments []*Argument // required

	groups []*Group

	envMap func(string) string

	usage *Usage
}

// Parse checks parsing errors and exits on errors
func (me *Parser) Parse() {
	if !me.Ok() {
		fmt.Println(me.Error())
		me.sh.Exit(1)
	}
}

func (me *Parser) SetShell(sh Shell) {
	me.sh = sh
	me.args = sh.Args()
}

func (me *Parser) Group(title, name string, items ...*Item) *Group {
	return me.group(title, name, me.Optional(name).String(""), items)
}

// Preface is the same as Usage().Preface
func (me *Parser) Preface(lines ...string) {
	me.usage.Preface(lines...)
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

// Usage returns the currently documented options for further
// documentation.
func (me *Parser) Usage() *Usage {
	return me.usage
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
		v:        me.parseMultiArg(name),
		required: true,
	}

	if len(arg.v) == 0 {
		arg.err = fmt.Errorf("missing %s", name)
	}
	me.arguments = append(me.arguments, arg)
	return arg
}

// Optional returns an optional  named argument.
func (me *Parser) Optional(name string) *Argument {
	arg := &Argument{
		name: name,
		v:    me.parseMultiArg(name),
	}
	me.arguments = append(me.arguments, arg)
	return arg
}

func (me *Parser) parseMultiArg(name string) []string {
	if isMulti(name) {
		return me.Args()
	}
	// parse one argument
	v := me.Argn(len(me.arguments))
	if v == "" {
		return nil
	}
	return []string{v}
}

type Argument struct {
	name     string
	v        []string
	err      error
	required bool
}

// String returns the value of this argument
func (me *Argument) String(def string) string {
	v := me.v
	if len(v) == 0 || v[0] == "" {
		return def
	}
	return v[0]
}

func isMulti(v string) bool {
	l := len(v)
	if l <= 3 {
		return false
	}
	end := v[l-3:]
	return end == "..."
}
