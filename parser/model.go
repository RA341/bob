package parser

import (
	"fmt"
	"strings"

	"github.com/RA341/bob/util"
)

type Arg struct {
	name       string
	defaultVal string
	// todo use enum
	argType  string
	required bool
}

type Task struct {
	cmd string
}

type Cmd struct {
	name string
	args map[string]Arg
	vars VarMap
	// commands to run
	tasks []Task
}

func (p Cmd) String() string {
	sb := strings.Builder{}

	sb.WriteString("Cmd: " + p.name + "\n")
	sb.WriteString("Args: " + fmt.Sprintf("%v", p.args) + "\n")
	sb.WriteString("Body: " + fmt.Sprintf("%v", p.tasks) + "\n")

	return sb.String()
}

type CommandMap map[string]Cmd

func (cm *CommandMap) GetCmdList() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for name, val := range *cm {
		var argList []string
		for v := range val.args {
			argList = append(argList, util.Yellow.Sprint(strings.TrimSpace(v)))
		}
		space := strings.Join(argList, ", ")

		sb.WriteString(util.Red.Sprint(" => ") + util.Green.Sprint(name) + "(" + space + ")" + "\n")
	}

	return sb.String()
}

type VarMap map[string]string

func (m *VarMap) Add(key, val string) {
	if *m == nil {
		*m = make(map[string]string)
	}
	(*m)[key] = val
}

type Bobfile struct {
	Version int
	Cmds    CommandMap
	Vars    VarMap
}

func (p Bobfile) String() string {
	return fmt.Sprintf(
		"Bobfile{Version: %v, Vars: %v}",
		p.Version,
		p.Vars,
	)
}
