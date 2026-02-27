package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/RA341/bob/util"
	"github.com/RA341/bob/vm"
)

type Bobfile struct {
	Version int
	Cmds    CommandMap

	Program vm.Program
}

func (b *Bobfile) String() string {
	return fmt.Sprintf(
		"Bobfile{Version: %v, Vars: %v}",
		b.Version,
		b.Program,
	)
}

func NewBobFileFromFile(bob *Bobfile, filepath string) error {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Could not read file", err)
	}

	return NewBobFileFromBytes(bob, contents)
}

func NewBobFileFromBytes(bob *Bobfile, fBytes []byte) error {
	contents := string(fBytes)
	contents = strings.TrimSpace(contents)

	lines := bob.splitAndStrip(contents)

	bob.Cmds = make(map[string]Cmd)

	//fmt.Println("Contents:")
	for _, line := range lines {
		//fmt.Printf("line %d: %s\n", i, line.content)
		if err := bob.parseLine(line); err != nil {
			return err
		}
	}

	return nil
}

const VersionPrefix = "version:"

func (b *Bobfile) parseLine(line CleanLine) error {
	if ok, err := b.parseVersion(line); err != nil {
		return err
	} else if ok {
		return nil
	}

	if ok, err := b.parseGlobalVar(line); err != nil {
		return err
	} else if ok {
		return nil
	}

	if b.parseCmd(line) {
		return nil
	}

	return nil
}

func (b *Bobfile) parseGlobalVar(line CleanLine) (bool, error) {
	if isInsideBraces(line.content, VarPrefix) {
		// this is a local var
		return false, nil
	}
	return b.parseVar(&line)
}

func (b *Bobfile) parseVar(line string) (ins []vm.Ins, err error) {
	if !isVar(line) {
		return false, nil
	}

	line = strings.TrimPrefix(line, VarPrefix)
	split := strings.Split(line, "=")
	if len(split) < 2 {
		return false, nil
	}

	// Left side of "=" is "key" or "key:type"
	argNameSplit := strings.Split(strings.TrimSpace(split[0]), ":")
	key := strings.TrimSpace(argNameSplit[0])

	typeName := vm.VTString
	if len(argNameSplit) > 1 {
		// has a type def
		typeName, err = vm.ValueTypeString(strings.TrimSpace(argNameSplit[1]))
		if err != nil {
			return false, fmt.Errorf(
				"invalid type name %q must be one of %v\nLine [%d]: %q",
				argNameSplit[1],
				vm.ValueTypeValues(),
				line,
			)
		}
	}

	// Right side of "=" is the value
	val := strings.TrimSpace(split[1])
	if val == "" {
		return false, fmt.Errorf("empty value: %q", line)
	}

	values := toPostfix(val)
	// parse expression
	if len(values) > 1 {
		var expr []vm.Ins

		for _, v := range values {
			space := strings.TrimSpace(v)
			switch {
			case space == "+":
				expr = append(expr, vm.O(vm.ADD))
			case isVar(space):
				expr = append(
					expr,
					vm.OVStr(vm.LOAD, strings.TrimPrefix(v, VarPrefix)),
				)
			}
		}

		b.Program.AddGlobalExpr(
			key,
			expr...,
		)
	} else {
		b.Program.AddGlobalVar(
			key,
			vm.NewValue(typeName, val),
		)
	}

	return true, nil
}

func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder

	for _, ch := range s {
		switch ch {
		case '+', '-', '*', '/', '(', ')':
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		case ' ':
			// skip
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}
	return tokens
}

func toPostfix(s string) []string {
	var output []string
	var ops []string

	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
	}

	tokens := tokenize(s)

	for _, tok := range tokens {
		switch tok {
		case "+", "-", "*", "/":
			for len(ops) > 0 {
				top := ops[len(ops)-1]
				if top != "(" && precedence[top] >= precedence[tok] {
					output = append(output, top)
					ops = ops[:len(ops)-1]
				} else {
					break
				}
			}
			ops = append(ops, tok)
		case "(":
			ops = append(ops, tok)
		case ")":
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				output = append(output, ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			if len(ops) > 0 {
				ops = ops[:len(ops)-1] // pop the "("
			}
		default:
			output = append(output, tok)
		}
	}

	for len(ops) > 0 {
		output = append(output, ops[len(ops)-1])
		ops = ops[:len(ops)-1]
	}

	return output
}

func (b *Bobfile) parseCmd(line CleanLine) bool {
	if line.content[0] >= '0' && line.content[0] <= '9' {
		log.Fatalf("Line %d starts with a number: %s:", line.lineNum, line.content)
	}

	curlyParen := 0
	roundParen := 0

	cmdName := ""
	args := ""
	argSet := false

	body := ""
	bodySet := false

	retVal := ""
	returnTypeSet := false
	util.UNUSED(retVal)

	acc := ""
	for _, l := range line.content {
		strL := string(l)

		switch strL {
		case "{":
			curlyParen++
		case "}":
			curlyParen--
		case "(":
			roundParen++
		case ")":
			roundParen--
		}

		// read till first ( -> gets the name
		if roundParen == 1 && cmdName == "" {
			if _, ok := b.Cmds[cmdName]; ok {
				log.Fatalf("Cmd redeclared %s", cmdName)
			}

			cmdName = strings.TrimSpace(acc)

			b.Cmds[cmdName] = Cmd{}
			acc = ""
			continue
		}

		// then read till first ) -> gets all args
		if strL == ")" && roundParen == 0 && !argSet {
			args = strings.TrimSpace(acc)
			argSet = true
			acc = ""
			continue
		}

		// then read till { -> potential return type
		if strL == "{" && curlyParen == 1 && !returnTypeSet {
			// this case is unused maybe useful in the future
			returnTypeSet = true
			retVal = strings.TrimSpace(acc)
			acc = ""
			continue
		}

		// then read till next } -> exec body
		if strL == "}" && curlyParen == 0 && !bodySet {
			bodySet = true
			body = strings.TrimSpace(acc)
			acc = ""
			continue
		}

		acc += strL
	}

	cmd, ok := b.Cmds[cmdName]
	if !ok {
		log.Fatalf("Could not find cmd %s, THIS SHOULD NEVER HAPPEN", cmdName)
	}

	cmd.args = b.convertArgs(args)
	cmd.tasks, cmd.vars = b.convertBody(body)
	cmd.name = cmdName
	b.Cmds[cmdName] = cmd

	return true
}

func (b *Bobfile) convertArgs(rawArgs string) map[string]Arg {
	args := make(map[string]Arg)
	if rawArgs == "" {
		return args
	}

	/*cases
	// valid
		<no args>
		user!
		user=default

	// invalid
		user!=default (cannot be required with default val)
	*/

	split := strings.Split(rawArgs, ",")
	for _, arg := range split {
		arg = cleanLine(arg)
		if arg == "" {
			continue
		}

		var ar Arg
		argSplit := strings.Split(arg, "=")
		// has default val
		if len(argSplit) > 1 {
			ar.defaultVal = cleanLine(argSplit[1])
		}

		argNameSplit := strings.Split(argSplit[0], ":")
		// does not have type info
		// user
		if len(argNameSplit) > 0 {
			argName := argNameSplit[0]
			if _, ok := args[argName]; ok {
				log.Fatalf("Argument '%s' already defined", argName)
			}

			ar.name = argName
		}

		if len(argNameSplit) > 1 {
			at := argNameSplit[1]
			if strings.HasSuffix(at, "!") {
				if ar.defaultVal != "" {
					log.Fatalf(
						"Argument '%s' has a default value %s, required args must not have default values",
						ar.name, ar.defaultVal,
					)
				}

				ar.required = true
			}

			ar.argType = strings.TrimSuffix(
				strings.TrimSpace(at),
				"!",
			)
		}

		args[ar.name] = ar
	}

	return args
}

func (b *Bobfile) convertBody(body string) ([]Task, VarMap) {
	var tasks []Task
	var varMap = VarMap{}

	var ins []vm.Ins

	for _, s := range strings.Split(body, "\n") {
		cleanCmd := cleanLine(s)
		if cleanCmd == "" || strings.HasPrefix(cleanCmd, "//") {
			continue
		}

		switch {
		case isFn(cleanCmd):

		case isVar(cleanCmd):
			varInitIns, err := b.parseVar(cleanCmd)
			if err != nil {
				// todo handle err
				log.Fatal(err)
			}

			ins = append(
				ins,
				varInitIns...,
			)
		}

		ts.cmd = cleanCmd
		tasks = append(tasks, ts)
	}

	return tasks, varMap
}

func isFn(cmd string) bool {

}

func (b *Bobfile) parseVersion(line CleanLine) (bool, error) {
	if !strings.HasPrefix(line.content, VersionPrefix) {
		return false, nil
	}

	splits := strings.Split(line.content, VersionPrefix)
	if len(splits) != 2 {
		return false, fmt.Errorf("could not parse version line, %s, %s", line.content, splits)
	}

	s := strings.TrimSpace(splits[1])
	atoi, err := strconv.Atoi(s)
	if err != nil {
		return false, fmt.Errorf("could not convert version to int, %s, %s, %v", line.content, s, err)
	}

	b.Version = atoi
	return true, nil
}

type CleanLine struct {
	lineNum int
	content string
}

func (b *Bobfile) splitAndStrip(contents string) []CleanLine {
	var cleanLines []CleanLine

	if !strings.HasSuffix(contents, "\n") {
		contents += "\n"
	}

	acc := ""
	curlyParen := 0
	roundParen := 0

	line := 0
	for _, ch := range contents {
		s := string(ch)

		switch s {
		case "{":
			curlyParen++
		case "}":
			curlyParen--
		case "(":
			roundParen++
		case ")":
			roundParen--
		}

		if s == "\n" {
			line++
			a := strings.TrimSpace(acc)

			if strings.HasPrefix(a, "//") {
				continue
			}

			// all scopes are closed
			if curlyParen == 0 && roundParen == 0 {
				if a != "" {
					//fmt.Printf("%d: %q\n", line, a)
					cleanLines = append(cleanLines, CleanLine{line, a})
				}
				acc = ""
				continue
			}
		}

		acc += s
	}

	if curlyParen != 0 || roundParen != 0 {
		log.Fatalf("Unclosed scope ")
	}

	return cleanLines
}
