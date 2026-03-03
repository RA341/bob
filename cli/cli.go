package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RA341/bob/cli/bob"
	"github.com/RA341/bob/cli/commands"
	"github.com/RA341/bob/cli/flag"
	"github.com/RA341/bob/util"
)

var ErrCmdHandled = errors.New("cmd handled")

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

func (c *Cli) Run(args []string) error {
	cleanArgs := args[1:]

	nextCmdIdx, err := c.parseFlags(cleanArgs)
	if err != nil {
		return fmt.Errorf("failed to parse commands: %w", err)
	}

	err = c.handleGlobalFlags()
	if err != nil {
		if errors.Is(err, ErrCmdHandled) {
			return nil
		}

		return err
	}

	cleanArgs = cleanArgs[nextCmdIdx:]

	return c.runSubCmd(cleanArgs)
}

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

func (c *Cli) handleGlobalFlags() error {
	if c.flags[0].IsSet {
		c.version()
		return ErrCmdHandled
	}

	if c.flags[1].IsSet {
		c.printHelp()
		return ErrCmdHandled
	}

	return nil
}

func (c *Cli) parseFlags(args []string) (subCmdIdx int, err error) {
	if len(args) < 1 {
		c.printHelp()
		return 0, nil
	}

	c.flags = append(
		c.flags,
		flag.HelpFlag("Print help"),
	)

	argParser, err := flag.ParseFlags(c.flags, args)
	if err != nil {
		return 0, err
	}

	return argParser.Current, nil
}

func (c *Cli) version() {
	fmt.Println("Bob version: dev")
}

func (c *Cli) printHelp() {
	fmt.Println(c.help())
}

func (c *Cli) help() string {
	var sb strings.Builder

	sb.WriteString(util.Cyan.Sprint("Bob command runner\n"))
	sb.WriteString(util.Yellow.Sprint("   Can we build it yes we can\n\n"))
	sb.WriteString(util.Cyan.Sprint("Subcommands\n"))
	for _, cmd := range c.cmds {
		sb.WriteString(Indent(
			cmd.Help(),
			util.Red.Sprint(" => "),
			"   ",
		))
	}

	sb.WriteString(util.Cyan.Sprint("To view help for a specific command:\n"))
	sb.WriteString(util.Magenta.Sprintf(
		"  %s <cmd> --help\n",
		filepath.Base(os.Args[0]),
	))

	return sb.String()
}

func Indent(input string, firstPad string, restPad string) string {
	// 1. Clean up the input to prevent trailing empty padded lines
	input = strings.TrimRight(input, "\n")
	if input == "" {
		return ""
	}

	// 2. Replace all internal newlines with the 'rest' padding
	// This ensures line 2, 3, etc. start with the correct indentation
	indented := strings.ReplaceAll(input, "\n", "\n"+restPad)

	// 3. Prepend the 'first' padding to the very beginning
	return fmt.Sprintf("%s%s\n", firstPad, indented)
}
