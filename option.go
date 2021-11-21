package cmdline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Option defines a command line option, ie. --username
type Option struct {
	args         []string // without command
	names        string
	defaultValue string
	quoteValue   bool // in usage output
	doc          []string
	err          error

	argIndex int // position in args for e.g. --username
	valIndex int // position for option value, same as argIndex if e.g. --i=1

	envMap func(string) string
}

// NewOption returns an option defined by a comma separated list of
// names and arguments to match against. Usually you would call
// Parser.Option(names) over this.
func NewOption(names string, args ...string) *Option {
	return &Option{names: names, args: args, argIndex: -1, valIndex: -1}
}

func (opt *Option) setDefault(def interface{}) {
	opt.defaultValue = fmt.Sprintf("%v", def)
}

// Doc sets the documentation lines for this option.
func (opt *Option) Doc(lines ...string) {
	opt.doc = lines
}

// Int same as IntOpt but does not return the Option.
func (opt *Option) Int(def int) int {
	v, _ := opt.IntOpt(def)
	return v
}

// IntOpt returns int value from the arguments or the given default value.
func (opt *Option) IntOpt(def int) (int, *Option) {
	opt.setDefault(def)
	v, err := opt.stringArg()
	if err != nil {
		opt.fail()
		return def, opt
	}
	if v == "" {
		return def, opt
	}
	iv, err := strconv.Atoi(v)
	if err != nil {
		opt.fail()
	}
	return iv, opt
}

// Uint8 same as uint8(Uint(...))
func (opt *Option) Uint8(def uint8) uint8 {
	return uint8(opt.Uint(uint64(def)))
}

// Uint16 same as uint16(Uint(...))
func (opt *Option) Uint16(def uint16) uint16 {
	return uint16(opt.Uint(uint64(def)))
}

// Uint32 same as uint32(Uint(...))
func (opt *Option) Uint32(def uint32) uint32 {
	return uint32(opt.Uint(uint64(def)))
}

// Uint same as UintOpt but does not return the Option
func (opt *Option) Uint(def uint64) uint64 {
	v, _ := opt.UintOpt(def)
	return v
}

// UintOpt returns an unsigned int option
func (opt *Option) UintOpt(def uint64) (uint64, *Option) {
	opt.setDefault(def)
	v, err := opt.stringArg()
	if err != nil {
		opt.fail()
		return def, opt
	}
	if v == "" {
		return def, opt
	}
	iv, err := strconv.ParseUint(v, 0, 64)
	if err != nil {
		opt.fail()
	}
	return iv, opt

}

// String same as StringOpt but does not return the Option.
func (opt *Option) String(def string) string {
	val, _ := opt.StringOpt(def)
	return val
}

// StringOpt returns string value from the arguments or the given default value.
func (opt *Option) StringOpt(def string) (string, *Option) {
	opt.setDefault(def)
	opt.quoteValue = true
	// todo , have to distinquish between option not found and value not found
	v, err := opt.stringArg()
	if err != nil {
		opt.fail()
		return def, opt
	}
	if isOption(v) {
		opt.fail()
	}
	if v == "" {
		return def, opt
	}
	return v, opt
}

func (opt *Option) stringArg() (string, error) {
	for i, arg := range opt.args {
		if opt.match(arg) {
			opt.argIndex = i
			opt.valIndex = i
			// NamedArg is -i=value
			eqIndex := strings.Index(arg, "=")
			if eqIndex > 0 {
				return arg[eqIndex+1:], nil
			}
			isLast := len(opt.args)-1 == i
			if isLast {
				opt.fail()
				return "", fmt.Errorf("missing value")
			}
			opt.valIndex = i + 1
			// NamedArg is -i
			return opt.args[i+1], nil
		}
	}
	return opt.envValue(), nil
}

// If last element in option names starts with $ expand it
func (opt *Option) envValue() string {
	names := opt.argNames()
	env := names[len(names)-1] // last element
	if env[0] != '$' {
		return ""
	}
	return os.Expand(env, opt.envMap)
}

func (opt *Option) argNames() []string {
	return strings.Split(strings.ReplaceAll(opt.names, " ", ""), ",")
}

func (opt *Option) match(arg string) bool {
	if !isOption(arg) {
		return false
	}
	names := opt.argNames()
	argName, _ := nameAndValue(arg)
	for _, name := range names {
		if name == argName {
			return true
		}
	}
	return false
}

func nameAndValue(arg string) (string, string) {
	parts := strings.Split(arg, "=")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func isOption(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

// Bool same as BoolOpt but does not return the Option.
func (opt *Option) Bool() bool {
	val, _ := opt.BoolOpt()
	return val
}

// BoolOpt returns bool value from the arguments or the given default
// value. The Option is returned for more configuration.
func (opt *Option) BoolOpt() (bool, *Option) {
	opt.setDefault("")
	v := opt.boolArg()
	return v, opt
}

func (opt *Option) boolArg() bool {
	for i, arg := range opt.args {
		if opt.match(arg) {
			opt.argIndex = i
			return true
		}
	}
	return false
}

func (opt *Option) fail() {
	opt.err = fmt.Errorf("Invalid option: %s", opt.names)
}

// Float64 returns float64
// value from the arguments or the given default value.
func (opt *Option) Float64(def float64) float64 {
	v, _ := opt.Float64Opt(def)
	return v
}

// Float64Opt returns float64 value from the arguments or the given
// default value.
func (opt *Option) Float64Opt(def float64) (float64, *Option) {
	opt.setDefault(def)
	v, err := opt.stringArg()
	if err != nil {
		opt.fail()
		return def, opt
	}
	if v == "" {
		return def, opt
	}
	iv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		opt.fail()
	}
	return iv, opt
}
