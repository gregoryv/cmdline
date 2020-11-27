package cmdline

type Item struct {
	Name   string
	Loader interface{}
}

// Load
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
