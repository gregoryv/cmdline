package cmdline

type Item interface {
	// Name must return one word
	Name() string

	// ExtraOptions for this action
	ExtraOptions(*Parser)
}
