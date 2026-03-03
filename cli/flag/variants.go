package flag

func IntFlag(name, help string) *Flag {
	// use 2 to indicate int otherwise it is considered bool
	return NewFlag(name, help, 2)
}

func StrFlag(name, help string) *Flag {
	return NewFlag(name, help, "")
}

func Bool(name, help string) *Flag {
	return NewFlag(name, help, false)
}

func NewFlag(name string, help string, val any) *Flag {
	return &Flag{
		Name: name,
		Val:  val,
		Help: help,
	}
}

// HelpFlag if help is empty it will set a default help message
func HelpFlag(help string) *Flag {
	if help == "" {
		help = "Print Help for this command"
	}
	return Bool("help", help)
}
