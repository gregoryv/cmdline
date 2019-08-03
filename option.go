package cmdline

import (
	"fmt"
	"strconv"
)

type Option struct {
	args         []string
	names        string
	defaultValue string
	doc          []string
	err          error
}

func (opt *Option) setDefault(def interface{}) {
	opt.defaultValue = fmt.Sprintf("%v", def)
}

func (opt *Option) Doc(lines ...string) {
	opt.doc = lines
}

func (opt *Option) Int(def int) int {
	v, _ := opt.IntOpt(def)
	return v
}

func (opt *Option) IntOpt(def int) (int, *Option) {
	opt.setDefault(def)
	v, _ := opt.string(opt.names)
	if v == "" {
		opt.fail()
		return def, opt
	}
	iv, err := strconv.Atoi(v)
	if err != nil {
		opt.fail()
	}
	return iv, opt
}

func (opt *Option) String(def string) string {
	val, _ := opt.StringOpt(def)
	return val
}

func (opt *Option) StringOpt(def string) (string, *Option) {
	opt.setDefault(def)
	v, _ := opt.string(opt.names)
	if v == "" {
		return def, opt
	}
	return v, opt
}

func (opt *Option) string(names string) (string, bool) {
	for i, arg := range opt.args {
		if namesMatch(names, arg) {
			isLast := len(opt.args)-1 == i
			if isLast {
				return "", true
			}
			return opt.args[i+1], true
		}
	}
	return "", false
}

func (opt *Option) Bool() bool {
	val, _ := opt.BoolOpt()
	return val
}

func (opt *Option) BoolOpt() (bool, *Option) {
	opt.setDefault(false)
	v, found := opt.string(opt.names)
	if v != "" && !isOption(v) {
		opt.fail()
		return false, opt
	}
	return found, opt
}

func (opt *Option) fail() {
	opt.err = fmt.Errorf("invalid %s", opt.names)
}
