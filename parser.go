// Package cmdline provides a way to parse command line arguments.
package cmdline

import (
	"fmt"
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
func (b *Basic) Parse() {
	b.defineHelp.Do(b.helpFlag)

	switch {
	case b.help:
		b.Usage().WriteTo(b.Parser.sh.Stdout())
		b.sh.Exit(0)

	case !b.Ok():
		fmt.Fprintln(b.sh.Stderr(), b.Error())
		fmt.Fprintln(b.sh.Stderr(), "Try -h or --help, for more information")
		b.sh.Exit(1)
	}
}

// Usage returns the usage for further documentation. If Parse method
// has not been called, it adds the help flag.
func (b *Basic) Usage() *Usage {
	b.defineHelp.Do(b.helpFlag)
	return b.Parser.Usage()
}

func (b *Basic) helpFlag() {
	b.help = b.Parser.Flag("-h, --help")
}

// ----------------------------------------

// DefaultShell is used by all new parsers.
var DefaultShell Shell = NewShellOS()

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
	arguments []*NamedArg // required

	groups []*Group

	envMap func(string) string

	usage *Usage
}

// Parse checks parsing errors and exits on errors
func (b *Parser) Parse() {
	if !b.Ok() {
		fmt.Fprintln(b.sh.Stderr(), b.Error())
		b.sh.Exit(1)
	}
}

func (b *Parser) SetShell(sh Shell) {
	b.sh = sh
	b.args = sh.Args()
}

func (b *Parser) Group(title, name string) *Group {
	return b.group(title, name, b.NamedArg(name).String(""))
}

// Preface is the same as Usage().Preface
func (b *Parser) Preface(lines ...string) {
	b.usage.Preface(lines...)
}

func (b *Parser) group(title, name, v string) *Group {
	grp := &Group{
		name:  name,
		args:  b.Args(),
		title: title,
		v:     v,
		items: make([]*Item, 0),
	}
	err := b.addGroup(grp)
	if err != nil {
		panic(fmt.Sprintf("duplicate group %q", title))
	}

	return grp
}

func (b *Parser) addGroup(grp *Group) error {
	for _, existing := range b.groups {
		if existing.Title() == grp.Title() {
			return fmt.Errorf("group %q already exists", grp.Title())
		}
	}
	b.groups = append(b.groups, grp)
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

func (b *Group) New(name string, any interface{}) *Item {
	item := &Item{
		Name:   name,
		Loader: any,
	}
	b.items = append(b.items, item)
	return item
}

// Selected returns the matching item. Defaults to the first in the group.
func (b *Group) Selected() interface{} {
	i := b.items[0]
	if b.v != "" {
		var found bool
		i, found = b.find(b.v)
		if !found {
			b.err = fmt.Errorf("invalid %s", b.name)
			return nil
		}
	}
	extra := NewParser()
	extra.args = append([]string{b.title}, b.args...)
	sel := i.Load(extra)
	b.err = extra.Error()
	return sel
}

func (b *Group) Title() string  { return b.title }
func (b *Group) Items() []*Item { return b.items }

// Find returns the named Item or nil if not found.
func (b *Group) find(name string) (*Item, bool) {
	for _, item := range b.items {
		if item.Name == name {
			return item, true
		}
	}
	return nil, false
}

// ----------------------------------------

// Ok returns true if no parsing error occured
func (b *Parser) Ok() bool {
	return b.Error() == nil
}

// Error returns first error of the given options.
func (b *Parser) Error() error {
	err := b.parseFailed()
	if err != nil {
		return err
	}
	if len(b.groups) == 0 { // as groups are selected with non option argument
		for _, arg := range b.Args() {
			if isOption(arg) {
				return fmt.Errorf("Unknown option: %v", arg)
			}
		}
	}
	return nil
}

func (b *Parser) parseFailed() error {
	var err error
	setErr := func(e error) {
		if err != nil {
			return
		}
		err = e
	}
	for _, opt := range b.options {
		setErr(opt.err)
	}
	for _, arg := range b.arguments {
		setErr(arg.err)
	}
	for _, grp := range b.groups {
		setErr(grp.err)
	}
	return err
}

// Option returns a new option with the given names.
// Names should be a comma separated string, e.g.
//
//	-n, --dry-run
//
// You can also include an environment variable to use as default
// value in the names, e.g.
//
//	-t, --token, $COGNITO_TOKEN
//
// The doclines are multiple documentation lines. The special docline
//
//	"hidden"
//
// means the values is masked when printed in the usage information.
func (b *Parser) Option(names string, doclines ...string) *Option {
	opt := NewOption(names, b.args[1:]...)
	opt.envMap = b.envMap
	opt.doc = make([]string, 0, len(doclines))
	for _, line := range doclines {
		if line == "hidden" {
			opt.hidden = true
			continue
		}
		opt.doc = append(opt.doc, line)

	}
	b.options = append(b.options, opt)
	return opt
}

// Flag is short for Option(name).Bool()
func (b *Parser) Flag(name string) bool {
	val, _ := b.Option(name).BoolOpt()
	return val
}

// Usage returns the currently documented options for further
// documentation.
func (b *Parser) Usage() *Usage {
	return b.usage
}

// Args returns arguments not matched by any of the options
func (b *Parser) Args() []string {
	rest := make([]string, 0)
	for i, arg := range b.args[1:] {
		//		fmt.Println("a:", arg)
		if !b.wasMatched(i) {
			rest = append(rest, arg)
		}
	}
	return rest
}

// Argn returns the n:th of remaining arguments starting at 0.
func (b *Parser) Argn(n int) string {
	rest := b.Args()
	if n < len(rest) {
		return rest[n]
	}
	return ""
}

func (b *Parser) wasMatched(i int) bool {
	for _, opt := range b.options {
		if opt.argIndex == i || opt.valIndex == i {
			return true
		}
	}
	return false
}

func (b *Parser) String() string {
	return fmt.Sprintf("Parser: %s", strings.Join(b.args, " "))
}

// NamedArg returns an named argument
func (b *Parser) NamedArg(name string) *NamedArg {
	arg := &NamedArg{
		name: name,
		v:    b.parseMultiArg(name),
	}
	b.arguments = append(b.arguments, arg)
	return arg
}

func (b *Parser) parseMultiArg(name string) []string {
	if isMulti(name) {
		return b.Args()
	}
	// parse one argument
	v := b.Argn(len(b.arguments))
	if v == "" {
		return nil
	}
	return []string{v}
}

// ----------------------------------------

type NamedArg struct {
	name     string
	v        []string
	err      error
	required bool
}

// String returns the value of this NamedArg or the given default
func (b *NamedArg) String(def string) string {
	b.required = (def == "")
	v := b.v

	if len(v) == 0 || v[0] == "" {
		return def
	}
	return v[0]
}

// Strings returns the values of this argument. If no default is given
// this NamedArg is considered required.
func (b *NamedArg) Strings(def ...string) []string {
	b.required = (len(def) == 0)
	switch {
	case len(b.v) == 0 && b.required:
		b.err = fmt.Errorf("missing %s", b.name)
	case len(b.v) == 0 && !b.required:
		return def
	}
	return b.v
}

func isMulti(v string) bool {
	l := len(v)
	if l <= 3 {
		return false
	}
	end := v[l-3:]
	return end == "..."
}
