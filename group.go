package cmdline

func NewGroup(title string, v ...Item) *Group {
	return &Group{
		title: title,
		items: v,
	}
}

type Group struct {
	title string
	items []Item
}

func (me *Group) Title() string { return me.title }
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
