package cmdline

func Example_runAction() {
	var (
		cli        = Parse("somecmd sayHi")
		actions, _ = cli.Group("ACTION", &Hi{})
		name       = cli.Required("ACTION").String()
	)
	action := actions.FindAction(name)
	action.ExtraOptions(cli)
	action.(Runnable).Run()
	// output:
	// Hi, stranger!
}

type Runnable interface {
	Run() error
}