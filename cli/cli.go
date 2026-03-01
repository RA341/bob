package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:     "greet",
		Usage:    "fight the loneliness!",
		Commands: []*cli.Command{},
		CommandNotFound: func(ctx context.Context, command *cli.Command, cmd string) {
			fmt.Println("Command not found:", cmd)
		},
	}
}
