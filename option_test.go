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
			got: New("", "-i", "1").Option("-i").Int(0),
			exp: 1,
		},
		{
			got: New("", "-i", "k").Option("-i").Int(0),
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
	}{}
	assert := asserter.New(t)
	for _, c := range cases {
		got := c.opt.match(c.arg)
		assert().Equals(got, c.exp)
	}
}
