package main

import (
	"fmt"
	"strings"
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

func (cm *CommandMap) GetCmdList() []string {
	var names []string
	for name := range *cm {
		names = append(names, name)
	}
	return names
}

type VarMap = map[string]string

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
