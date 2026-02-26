package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Runner struct {
	bf        *Bobfile
	finalArgs []string
	argLen    int
	ctx       context.Context
}

func NewRunner(ctx context.Context, bf *Bobfile) {
	r := Runner{
		bf:  bf,
		ctx: ctx,
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: bob <command>")
		r.printAvailCmd()
		return
	}

	r.finalArgs = os.Args[1:]
	r.argLen = len(r.finalArgs)

	//fmt.Println("Args: ", r.finalArgs)

	r.Execute()
}

func (r *Runner) Execute() {
	if r.argLen > 0 {
		task := r.finalArgs[0]
		val, ok := r.bf.Cmds[task]
		if !ok {
			fmt.Printf("Unknown cmd: %s\n", task)
			r.printAvailCmd()
		}

		r.runSubCmd(&val)
	}
}

func (r *Runner) printAvailCmd() {
	cl := r.bf.Cmds.GetCmdList()
	fmt.Println("Available commands: ", cl)
}

func (r *Runner) runSubCmd(c *Cmd) {
	for _, t := range c.tasks {
		if r.runBuiltin(&t, c) {
			continue
		}

		r.runShell(&t, c)
	}
}

func (r *Runner) runShell(t *Task, cmd *Cmd) {
	t.cmd = r.replaceArgIfAny(t.cmd, cmd)
	fmt.Println("Executing shell [", t.cmd, "]")

	execStr := strings.Split(t.cmd, " ")
	ex := exec.CommandContext(r.ctx, execStr[0], execStr[1:]...)
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	ex.Stdin = os.Stdin

	err := ex.Run()
	if err != nil {
		log.Fatal("Error running command", err)
	}
}

// parses and runs a builtin if exists
func (r *Runner) runBuiltin(t *Task, cmd *Cmd) bool {
	acc := ""

	name := ""
	args := ""
	var execFn BuiltInFunc

	for _, sd := range t.cmd {
		sdStr := string(sd)
		if sdStr == "(" {
			name = acc

			var ok bool
			execFn, ok = builtInMaps[name]
			if !ok {
				return false
			}

			acc = ""
			continue
		}

		if sdStr == ")" {
			args = acc
		}

		acc += sdStr
	}

	if execFn == nil {
		return false
	}

	args = r.replaceArgIfAny(args, cmd)
	fmt.Println("Executing builtin [", name, args, "]")

	err := execFn(strings.Split(args, ",")...)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func (r *Runner) replaceArgIfAny(args string, cmd *Cmd) string {
	re := regexp.MustCompile(`\$\{([^}]+)}`)
	matches := re.FindAllStringSubmatch(args, -1)

	replaceFn := func(argName, actualVal string) string {
		replaceStr := fmt.Sprintf("${%s}", argName)
		args = strings.ReplaceAll(args, replaceStr, actualVal)

		//fmt.Printf("%s = %s\n", argName, replaceStr)
		return replaceStr
	}

	for _, x := range matches {
		// get value inside ${}
		argName := x[1]

		argVal, ok := cmd.args[argName]
		if ok {
			val := r.getArgVal(&argVal, cmd)
			replaceFn(argName, val)
			continue
		}

		val, ok := r.bf.Vars[argName]
		if ok {
			replaceFn(argName, val)
			continue
		}
	}

	return args
}

func (r *Runner) getArgVal(a *Arg, cmd *Cmd) string {
	argCmds := r.finalArgs[1:]

	for _, ac := range argCmds {
		splt := strings.Split(ac, "=")
		if len(splt) == 2 && splt[0] == a.name {
			return splt[1]
		}
	}

	if a.required {
		example := fmt.Sprintf("%s %s %s=value", os.Args[0], cmd.name, a.name)
		log.Fatal("Argument ", a.name, " is required, pass it like so ", example)
	}

	return a.defaultVal
}
