package cmdline

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_flags(t *testing.T) {
	cases := []struct {
		opt *Option
		exp bool
	}{
		{
			opt: New("", "-n").Option("-n"),
			exp: true,
		},
		{
			opt: New("", "-n", "-name").Option("-n"),
			exp: true,
		},
		{
			opt: New("", "-n", "val").Option("-n"),
			exp: false,
		},
	}

	assert := asserter.New(t)
	for _, c := range cases {
		got := c.opt.Bool()
		assert().Equals(got, c.exp)
	}
}
