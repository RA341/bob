package cli

import (
	"fmt"
	"log"
	"os"
)

type Command interface {
	Name() string
	Help() string
	Run() error
}
type Cli struct {
	cmds []Command
}

func Run() {
	cli := Cli{
		cmds: []Command{
			&CmdInit{},
		},
	}

	if len(os.Args) < 2 {
		cli.PrintHelp()
		return
	}

	subCommand := os.Args[1]

	for _, cmd := range cli.cmds {
		if cmd.Name() == subCommand {
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
				return
			}
			return
		}
	}

	RunBobFile(subCommand)
}

func (c *Cli) PrintHelp() {
	fmt.Println("Bob command runner")
	fmt.Println("   Can we build it yes we can")
	fmt.Println()

	fmt.Println("Subcommands:")
	for _, cmd := range c.cmds {
		name := fmt.Sprintf("   %s: %s", cmd.Name(), cmd.Help())
		fmt.Println(name)
	}

}
