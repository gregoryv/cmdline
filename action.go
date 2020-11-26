package cmdline

type Action interface {
	// Name must return one word
	Name() string

	// ExtraOptions for this action
	ExtraOptions(*CommandLine)
}
