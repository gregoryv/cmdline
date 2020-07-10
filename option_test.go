package cmdline

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestOption_Uint(t *testing.T) {
	assert := asserter.New(t)
	assert().Equals(NewOption("-o", "-o", "1").Uint(2), uint64(1))
	assert().Equals(NewOption("-o", "-o", "7").Uint(2), uint64(7))
	assert().Equals(NewOption("-o", "-o", "-9").Uint(2), uint64(0))
	assert(NewOption("-o", "-o").Uint(1) == 0).Error("missing value")
	assert(NewOption("-o", "-o", "string").Uint(1) == 0).Error("bad value")

	assert(NewOption("-o", "-o", "1").Int(2) == 1).Fail()
	assert(NewOption("-o", "-o").Int(1) == 0).Error("missing value")
	assert(NewOption("-o", "-o", "string").Int(1) == 0).Error("bad value")

	assert(NewOption("-n", "-n").Bool())
	assert(NewOption("-n", "-n", "-name").Bool())
	assert(!NewOption("-n", "-n", "string").Bool())
}
