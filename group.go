package cmdline

func NewGroup(name string, v ...Action) *Group {
	return &Group{
		name:    name,
		actions: v,
	}
}

type Group struct {
	name    string
	actions []Action
}

func (me *Group) Name() string { return me.name }

// FindAction returns the named action or nil if not found.
func (me *Group) FindAction(name string) Action {
	for _, a := range me.actions {
		if a.Name() == name {
			return a
		}
	}
	return nil
}
