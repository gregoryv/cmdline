package cmdline

import "testing"

func TestGroup_FindAction_returns_nil_on_undefined_action(t *testing.T) {
	grp := NewGroup("empty")
	if grp.FindAction("any") != nil {
		t.Fail()
	}
}
