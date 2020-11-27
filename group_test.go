package cmdline

import "testing"

func TestGroup_Find_returns_nil_on_undefined_item(t *testing.T) {
	grp := NewGroup("empty")
	if grp.Find("any") != nil {
		t.Fail()
	}
}
