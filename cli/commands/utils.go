package commands

import (
	"strings"

	"github.com/RA341/bob/cli/flag"
	"github.com/RA341/bob/util"
)

func PrintHelpSubCmd(name string, flags []*flag.Flag, helpText string) string {
	var sb strings.Builder
	sb.WriteString(util.Cyan.Sprintf("%s:\n", name))
	sb.WriteString(util.Green.Sprint("  Usage: " + helpText + "\n"))
	sb.WriteString(util.Yellow.Sprint("  " + "Flags:" + "\n"))

	for _, s := range flags {
		sb.WriteString(util.Yellow.Sprint("    --"+s.Name) + util.Blue.Sprint(" "+s.Help+"\n"))
	}

	return sb.String()
}
