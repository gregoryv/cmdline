package cmdline

import "testing"

func TestGroup_Find_returns_nil_on_undefined_item(t *testing.T) {
	grp := NewGroup("empty")
	if _, found := grp.Find("any"); found {
		t.Fail()
	}
}
