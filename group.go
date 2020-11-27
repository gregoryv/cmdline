package cmdline

func NewGroup(name string, v ...Item) *Group {
	return &Group{
		name:  name,
		items: v,
	}
}

type Group struct {
	name  string
	items []Item
}

func (me *Group) Name() string  { return me.name }
func (me *Group) Items() []Item { return me.items }

// Find returns the named Item or nil if not found.
func (me *Group) Find(name string) (Item, bool) {
	for _, a := range me.items {
		if a.Name() == name {
			return a, true
		}
	}
	return nil, false
}
