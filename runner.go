package main

import (
	"fmt"
	"log"
	"os"
)

func Execute(bf *Bobfile) {
	if len(os.Args) < 1 {
		fmt.Println("Usage: bob <command>")
		cl := bf.Cmds.GetCmdList()
		fmt.Println("Available commands: ", cl)
		return
	}

	realArgs := os.Args[1:]
	argLen := len(realArgs)

	if argLen > 0 {
		task := realArgs[0]
		val, ok := bf.Cmds[task]
		if !ok {
			fmt.Printf("Unknown cmd: %s\n", task)
			cl := bf.Cmds.GetCmdList()
			fmt.Println("Available commands: ", cl)
		}

		runSubCmd(bf, &val)
	}
}

type BuiltInFunc func(vals ...interface{}) error

var builtInMaps = map[string]BuiltInFunc{
	"print": func(vals ...interface{}) error {
		fmt.Println(vals...)
		return nil
	},
}

func runSubCmd(bf *Bobfile, c *Cmd) {
	for _, t := range c.tasks {
		if checkAndRunBuiltin(bf, &t) {
			continue
		}

		runShell(bf, &t)
	}
}

func runShell(bf *Bobfile, t *Task) {
	fmt.Println("Running in Shell", t.cmd)
}

// parses and runs a builtin if exists
func checkAndRunBuiltin(bf *Bobfile, t *Task) bool {
	acc := ""

	for _, sd := range t.cmd {
		sdStr := string(sd)
		if sdStr == "(" {
			val, ok := builtInMaps[acc]
			if !ok {
				return false
			}

			fmt.Println("Executing:", acc)
			err := val("Some val")
			if err != nil {
				log.Fatal(err)
			}

			acc = ""
			continue
		}

		acc += sdStr
	}

	return false
}
