package cli

import (
	"fmt"
	"log"
	"os"
)

type Command interface {
	Name() string
	Help() string
	Run(args ...string) error
}

type Cli struct {
	cmds      []Command
	flags     []*Flag
	subCmdIdx int
}

func Run() {
	versionFlag := BoolFlag("version", "print version")
	cli := Cli{
		flags: []*Flag{
			versionFlag,
		},
		cmds: []Command{
			&CmdInit{},
		},
	}

	cleanArgs := os.Args[1:]

	err := cli.ParseFlags(cleanArgs)
	if err != nil {
		log.Fatal("Could not parse flags", err)
	}

	if versionFlag.isSet {
		cli.Version()
		return
	}

	subCommand := cleanArgs[cli.subCmdIdx]
	subCommandArgs := cleanArgs[cli.subCmdIdx+1:]

	for _, cmd := range cli.cmds {
		if cmd.Name() == subCommand {
			err := cmd.Run(subCommandArgs...)
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

func (c *Cli) ParseFlags(args []string) error {
	if len(args) < 1 {
		c.PrintHelp()
		return nil
	}

	argParser := FlagParser{
		flags: c.flags,
	}
	err := argParser.parse(args)
	if err != nil {
		return err
	}

	c.flags = argParser.flags
	c.subCmdIdx = argParser.current

	return nil
}

func (c *Cli) Version() {
	fmt.Println("Bob version: dev")
}
