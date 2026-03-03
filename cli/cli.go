package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/RA341/bob/cli/bob"
	"github.com/RA341/bob/cli/commands"
	"github.com/RA341/bob/cli/flag"
)

type Command interface {
	Name() string
	Help() string
	Run() error
	LoadFlags(args []string) (err error)
}

type Cli struct {
	cmds      []Command
	flags     []*flag.Flag
	subCmdIdx int
}

func NewApp() *Cli {
	return &Cli{
		flags: []*flag.Flag{
			flag.Bool("version", "print version"),
		},
		cmds: []Command{
			new(commands.Init),
		},
	}
}

func (c *Cli) Run() error {
	cleanArgs := os.Args[1:]

	nextCmdIdx, err := c.ParseFlags(cleanArgs)
	if err != nil {
		return fmt.Errorf("failed to parse commands: %w", err)
	}

	err = c.handleGlobalFlags()
	if err != nil && !errors.Is(err, ErrCmdHandled) {
		return err
	}

	cleanArgs = cleanArgs[nextCmdIdx:]

	return c.runSubCmd(cleanArgs)
}

func (c *Cli) handleGlobalFlags() error {
	if c.flags[0].IsSet {
		c.Version()
		return ErrCmdHandled
	}

	return nil
}

var ErrCmdHandled = errors.New("cmd handled")

func (c *Cli) runSubCmd(cleanArgs []string) (err error) {
	subCommand := cleanArgs[c.subCmdIdx]
	subCommandArgs := cleanArgs[c.subCmdIdx+1:]

	// todo add debug printer
	//fmt.Println("Subcommand", subCommand)
	//fmt.Println("Args", subCommandArgs)

	for _, cmd := range c.cmds {
		if cmd.Name() != subCommand {
			continue
		}

		if err = cmd.LoadFlags(subCommandArgs); err != nil {
			return err
		}

		return cmd.Run()
	}

	return bob.Run(subCommand, subCommandArgs)
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

func (c *Cli) ParseFlags(args []string) (subCmdIdx int, err error) {
	if len(args) < 1 {
		c.PrintHelp()
		return 0, nil
	}

	argParser, err := flag.ParseFlags(c.flags, args)
	if err != nil {
		return 0, err
	}

	return argParser.Current, nil
}

func (c *Cli) Version() {
	fmt.Println("Bob version: dev")
}
