package cmdline

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_flags(t *testing.T) {
	cases := []struct {
		got, exp bool
	}{
		{New("", "-n").Option("-n").Bool(), true},
		{New("", "-n", "-name").Option("-n").Bool(), true},
		{New("", "-n", "val").Option("-n").Bool(), false},
	}

	assert := asserter.New(t)
	for _, c := range cases {
		assert().Equals(c.got, c.exp)
	}
}

func Test_int_options(t *testing.T) {
	cases := []struct {
		got int
		exp int
	}{
		{
			got: NewOption("-i", "-i", "1").Int(0),
			exp: 1,
		},
		{
			got: NewOption("-i", "-i=1").Int(0),
			exp: 1,
		},
		{
			got: NewOption("-i", "-i", "k").Int(0),
			exp: 0,
		},
	}
	assert := asserter.New(t)
	for _, c := range cases {
		assert().Equals(c.got, c.exp)
	}
}

func Test_namesMatch(t *testing.T) {
	cases := []struct {
		opt *Option
		arg string
		exp bool
	}{
		{NewOption("-v"), "-v", true},
		{NewOption("-i"), "-i=1", true},
	}
	assert := asserter.New(t)
	for _, c := range cases {
		got := c.opt.match(c.arg)
		assert().Equals(got, c.exp)
	}
}
