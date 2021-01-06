package cmdline

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func Test_default_string_option(t *testing.T) {
	cli := Parse("mycmd")
	got := cli.Option("-b").String("default")
	if got != "default" {
		t.Error("unexpected:", got)
	}
}

func Test_missing_string_option(t *testing.T) {
	cli := Parse("mycmd -a -b=1")
	got := cli.Option("-a").String("")
	cli.Option("-b").Int(0)
	if cli.Ok() {
		t.Error("should fail:", got)
	}
}

func Test_quoted_string_option(t *testing.T) {
	cli := Parse(`mycmd -a "-b=1"`)
	got := cli.Option("-a").String("")
	if !cli.Ok() {
		t.Error(cli.Error(), got)
	}
}

func TestOption(t *testing.T) {
	assert := asserter.New(t)

	var (
		cli     = Parse("countstars -verbose -min 1 -filter alien -last")
		min     = cli.Option("-min")
		verbose = cli.Option("-verbose")
		max     = cli.Option("-max")
		last    = cli.Option("-last")
		ux      = cli.Option("-ux")
	)

	assert().Equals(ux.Uint8(8), uint8(8))
	assert().Equals(ux.Uint16(16), uint16(16))
	assert().Equals(ux.Uint32(32), uint32(32))

	assert().Equals(min.Uint(2), uint64(1))
	assert().Equals(verbose.Uint(2), uint64(0))
	assert().Equals(max.Uint(2), uint64(2))
	//assert().Equals(last.Uint(2), uint64(0))

	assert().Equals(min.Int(2), 1)
	assert().Equals(verbose.Int(2), 0)
	assert().Equals(max.Int(2), 2)

	assert().Equals(min.Bool(), true)
	assert().Equals(verbose.Bool(), true)
	assert().Equals(max.Bool(), false)

	assert().Equals(min.String("s"), "1")
	assert().Equals(verbose.String("s"), "-min")
	assert(verbose.err != nil).Error(verbose.err)

	assert().Equals(max.String("s"), "s")

	// fixme
	//assert().Equals(last.String("s"), "")
	//assert(last.err != nil).Error(last.err)
	_ = last
}

func TestOption_String(t *testing.T) {
	got := NewOption("-n", "-n", "john").String("doe")
	assert := asserter.New(t)
	assert().Equals(got, "john")
}
