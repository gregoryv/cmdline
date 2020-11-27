package cmdline

type Item string

func (me Item) Name() string { return string(me) }
