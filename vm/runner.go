package vm

//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//	"os/exec"
//	"regexp"
//	"strings"
//
//	"github.com/fatih/color"
//)
//
//type RuntimeEnv struct {
//	workingDir         string
//	originalWorkingDir string
//}
//
//type Runner struct {
//	bf        *Bobfile
//	rte       *RuntimeEnv
//	finalArgs []string
//	argLen    int
//	ctx       context.Context
//}
//
//func NewRunner(ctx context.Context, bf *Bobfile, workingDir string) {
//	r := Runner{
//		bf:  bf,
//		ctx: ctx,
//		rte: &RuntimeEnv{
//			workingDir:         workingDir,
//			originalWorkingDir: workingDir,
//		},
//	}
//
//	if len(os.Args) < 2 {
//		blue.PrintlnFunc()("Usage: bob <command>")
//		r.printAvailCmd()
//		return
//	}
//
//	r.finalArgs = os.Args[1:]
//	r.argLen = len(r.finalArgs)
//
//	//fmt.Println("Args: ", r.finalArgs)
//
//	r.Execute()
//}
//
//func (r *Runner) Execute() {
//	if r.argLen > 0 {
//		task := r.finalArgs[0]
//		val, ok := r.bf.Cmds[task]
//		if !ok {
//			red.PrintfFunc()("Unknown cmd: %s\n", task)
//			r.printAvailCmd()
//		}
//
//		r.runSubCmd(&val)
//	}
//}
//
//func (r *Runner) printAvailCmd() {
//	cl := r.bf.Cmds.GetCmdList()
//	fmt.Println(blue.Sprint("Available commands: "), cl)
//}
//
//func (r *Runner) runSubCmd(c *Cmd) {
//	for _, t := range c.tasks {
//		if r.runBuiltin(&t, c) {
//			continue
//		}
//
//		r.runShell(&t, c)
//	}
//}
//
//var blue = color.New(color.FgCyan)
//
//func (r *Runner) runShell(t *Task, cmd *Cmd) {
//	t.cmd = r.replaceArgIfAny(t.cmd, cmd)
//
//	_, _ = blue.Println("[", t.cmd, "]")
//
//	execStr := strings.Split(t.cmd, " ")
//	ex := exec.CommandContext(r.ctx, execStr[0], execStr[1:]...)
//	sprintf := blue.Sprintf("=> ")
//
//	ex.Dir = r.rte.workingDir
//	ex.Stdout = NewPrefixWriter(os.Stdout, sprintf)
//	ex.Stderr = NewPrefixWriter(os.Stderr, sprintf)
//	ex.Stdin = os.Stdin
//
//	err := ex.Run()
//	if err != nil {
//		fmt.Println(sprintf, red.Sprintf("err: %s", err.Error()))
//	}
//}
//
//var green = color.New(color.FgHiGreen)
//
//// parses and runs a builtin if exists
//func (r *Runner) runBuiltin(t *Task, cmd *Cmd) bool {
//	name, args, execFn := r.parseBuiltIn(t)
//	if execFn == nil {
//		return false
//	}
//
//	args = r.replaceArgIfAny(args, cmd)
//	fmt.Println(
//		green.Sprint("[", name),
//		yellow.Sprint(args),
//		green.Sprint("]"),
//	)
//
//	err := execFn(r.rte, strings.Split(args, ",")...)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return true
//}
//
//func (r *Runner) parseBuiltIn(t *Task) (name string, args string, runner BuiltInFunc) {
//	var execFn BuiltInFunc
//
//	acc := ""
//	for _, sd := range t.cmd {
//		sdStr := string(sd)
//		if sdStr == "(" {
//			name = acc
//
//			var ok bool
//
//			execFn, ok = builtInMaps[name]
//			if !ok {
//				return "", "", nil
//			}
//
//			//cmd, ok := r.bf.Cmds[name]
//			//if ok {
//			//	execFn = func(env *RuntimeEnv, vals ...string) error {
//			//
//			//
//			//		r.runBuiltin(t, &cmd)
//			//
//			//		return nil
//			//	}
//			//}
//
//			acc = ""
//			continue
//		}
//
//		if sdStr == ")" {
//			args = acc
//		}
//
//		acc += sdStr
//	}
//
//	return name, args, execFn
//}
//
//func (r *Runner) replaceArgIfAny(args string, cmd *Cmd) string {
//	re := regexp.MustCompile(`\$\{([^}]+)}`)
//	matches := re.FindAllStringSubmatch(args, -1)
//
//	replaceFn := func(argName, actualVal string) string {
//		replaceStr := fmt.Sprintf("${%s}", argName)
//		args = strings.ReplaceAll(args, replaceStr, actualVal)
//
//		//fmt.Printf("%s = %s\n", argName, replaceStr)
//		return replaceStr
//	}
//
//	for _, x := range matches {
//		// get value inside ${}
//		argName := x[1]
//		if argName == "" {
//			continue
//		}
//
//		argVal, ok := cmd.args[argName]
//		if ok {
//			val := r.getArgVal(&argVal, cmd)
//			replaceFn(argName, val)
//			continue
//		}
//
//		val, ok := cmd.vars[argName]
//		if ok {
//			replaceFn(argName, val)
//			continue
//		}
//
//		val, ok = r.bf.Vars[argName]
//		if ok {
//			replaceFn(argName, val)
//			continue
//		}
//	}
//
//	return args
//}
//
//func (r *Runner) getArgVal(a *Arg, cmd *Cmd) string {
//	argCmds := r.finalArgs[1:]
//
//	for _, ac := range argCmds {
//		splt := strings.Split(ac, "=")
//		if len(splt) == 2 && splt[0] == a.name {
//			return splt[1]
//		}
//	}
//
//	if a.required {
//		example := fmt.Sprintf("%s %s %s=value", os.Args[0], cmd.name, a.name)
//		log.Fatal("Argument ", a.name, " is required, pass it like so ", example)
//	}
//
//	return a.defaultVal
//}
