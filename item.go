package cmdline

type Item struct {
	Name   string
	Loader interface{}
}

// Load returns the item with extra options if it implements
// WithExtraOptions interface.
func (me *Item) Load(p *Parser) interface{} {
	switch l := me.Loader.(type) {
	case func(*Parser) interface{}:
		return l(p)
	case WithExtraOptions:
		l.ExtraOptions(p)
		return l
	default:
		return l
	}
}

type WithExtraOptions interface {
	// ExtraOptions is used to parse extra options for a grouped item
	ExtraOptions(*Parser)
}
