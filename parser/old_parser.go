package parser

//import (
//	"fmt"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//
//	"github.com/RA341/bob/util"
//	"github.com/RA341/bob/vm"
//)
//
//type Bobfile struct {
//	Version int
//	Cmds    CommandMap
//
//	Program vm.Program
//}
//
//func (b *Bobfile) String() string {
//	return fmt.Sprintf(
//		"Bobfile{Version: %v, Vars: %v}",
//		b.Version,
//		b.Program,
//	)
//}
//
//func NewBobFileFromFile(bob *Bobfile, filepath string) error {
//	contents, err := os.ReadFile(filepath)
//	if err != nil {
//		log.Fatal("Could not read file", err)
//	}
//
//	return NewBobFileFromBytes(bob, contents)
//}
//
//func NewBobFileFromBytes(bob *Bobfile, fBytes []byte) error {
//	contents := string(fBytes)
//	contents = strings.TrimSpace(contents)
//
//	lines := splitAndStrip(contents)
//
//	bob.Cmds = make(map[string]Cmd)
//
//	//fmt.Println("Contents:")
//	for _, line := range lines {
//		//fmt.Printf("line %d: %s\n", i, line.content)
//		if err := bob.parseLine(line); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//const VersionPrefix = "version:"
//
//func (b *Bobfile) parseLine(line CleanLine) error {
//	if ok, err := b.parseVersion(line); err != nil {
//		return err
//	} else if ok {
//		return nil
//	}
//
//	if ok, err := b.parseGlobalVar(line); err != nil {
//		return err
//	} else if ok {
//		return nil
//	}
//
//	if b.parseCmd(line) {
//		return nil
//	}
//
//	return nil
//}
//
//func (b *Bobfile) parseGlobalVar(line CleanLine) (bool, error) {
//	//if isInsideBraces(line.content, VarPrefix) {
//	//	// this is a local var
//	//	return false, nil
//	//}
//
//	if !isVar(line.content) {
//		return false, nil
//	}
//
//	ins, err := b.parseVar(line.content)
//	b.Program.Add(ins...)
//
//	return true, err
//}
//
//func (b *Bobfile) parseVar(line string) (ins []vm.Ins, err error) {
//	if !isVar(line) {
//		return nil, nil
//	}
//
//	line = strings.TrimPrefix(line, VarPrefix)
//	split := strings.Split(line, "=")
//	if len(split) < 2 {
//		return nil, nil
//	}
//
//	// Left side of "=" is "key" or "key:type"
//	argNameSplit := strings.Split(strings.TrimSpace(split[0]), ":")
//	key := strings.TrimSpace(argNameSplit[0])
//
//	typeName := vm.VTString
//	if len(argNameSplit) > 1 {
//		// has a type def
//		typeName, err = vm.ValueTypeString(strings.TrimSpace(argNameSplit[1]))
//		if err != nil {
//			return nil, fmt.Errorf(
//				"invalid type name %q must be one of %v\n: %q",
//				argNameSplit[1],
//				vm.ValueTypeValues(),
//				line,
//			)
//		}
//	}
//
//	// Right side of "=" is the value
//	val := strings.TrimSpace(split[1])
//	if val == "" {
//		return nil, fmt.Errorf("empty value: %q", line)
//	}
//
//	values := toPostfix(val)
//	// parse expression
//	if len(values) > 1 {
//		var expr []vm.Ins
//
//		for _, v := range values {
//			space := strings.TrimSpace(v)
//			switch {
//			case space == "+":
//				expr = append(expr, vm.O(vm.ADD))
//			case isVar(space):
//				expr = append(
//					expr,
//					vm.OVStr(vm.LOAD, strings.TrimPrefix(v, VarPrefix)),
//				)
//			}
//		}
//
//		return vm.AddExpr(
//			key,
//			expr...,
//		), nil
//	}
//
//	return vm.AddVar(key, vm.NewValue(typeName, val)), nil
//}
//
//func (b *Bobfile) parseCmd(line CleanLine) bool {
//	if line.content[0] >= '0' && line.content[0] <= '9' {
//		log.Fatalf("Line %d starts with a number: %s:", line.lineNum, line.content)
//	}
//
//	curlyParen := 0
//	roundParen := 0
//
//	cmdName := ""
//	args := ""
//	util.UNUSED(args)
//	argSet := false
//
//	body := ""
//	bodySet := false
//
//	retVal := ""
//	returnTypeSet := false
//	util.UNUSED(retVal)
//
//	acc := ""
//	for _, l := range line.content {
//		strL := string(l)
//
//		switch strL {
//		case "{":
//			curlyParen++
//		case "}":
//			curlyParen--
//		case "(":
//			roundParen++
//		case ")":
//			roundParen--
//		}
//
//		// read till first ( -> gets the name
//		if roundParen == 1 && cmdName == "" {
//			if _, ok := b.Cmds[cmdName]; ok {
//				log.Fatalf("Cmd redeclared %s", cmdName)
//			}
//
//			cmdName = strings.TrimSpace(acc)
//
//			b.Cmds[cmdName] = Cmd{}
//			acc = ""
//			continue
//		}
//
//		// then read till first ) -> gets all args
//		if strL == ")" && roundParen == 0 && !argSet {
//			args = strings.TrimSpace(acc)
//			argSet = true
//			acc = ""
//			continue
//		}
//
//		// then read till { -> potential return type
//		if strL == "{" && curlyParen == 1 && !returnTypeSet {
//			// this case is unused maybe useful in the future
//			returnTypeSet = true
//			retVal = strings.TrimSpace(acc)
//			acc = ""
//			continue
//		}
//
//		// then read till next } -> exec body
//		if strL == "}" && curlyParen == 0 && !bodySet {
//			bodySet = true
//			body = strings.TrimSpace(acc)
//			acc = ""
//			continue
//		}
//
//		acc += strL
//	}
//
//	cmd, ok := b.Cmds[cmdName]
//	if !ok {
//		log.Fatalf("Could not find cmd %s, THIS SHOULD NEVER HAPPEN", cmdName)
//	}
//
//	// todo args
//	//cmd.args = b.convertArgs(args)
//	bodyIns := b.convertBody(body)
//
//	b.Program.Add(vm.OVStr(vm.LABEL, cmdName))
//	b.Program.Add(bodyIns...)
//
//	b.Cmds[cmdName] = cmd
//
//	return true
//}
//
//func (b *Bobfile) convertArgs(rawArgs string) map[string]Arg {
//	args := make(map[string]Arg)
//	if rawArgs == "" {
//		return args
//	}
//
//	/*cases
//	// valid
//		<no args>
//		user!
//		user=default
//
//	// invalid
//		user!=default (cannot be required with default val)
//	*/
//
//	split := strings.Split(rawArgs, ",")
//	for _, arg := range split {
//		arg = cleanLine(arg)
//		if arg == "" {
//			continue
//		}
//
//		var ar Arg
//		argSplit := strings.Split(arg, "=")
//		// has default val
//		if len(argSplit) > 1 {
//			ar.defaultVal = cleanLine(argSplit[1])
//		}
//
//		argNameSplit := strings.Split(argSplit[0], ":")
//		// does not have type info
//		// user
//		if len(argNameSplit) > 0 {
//			argName := argNameSplit[0]
//			if _, ok := args[argName]; ok {
//				log.Fatalf("Argument '%s' already defined", argName)
//			}
//
//			ar.name = argName
//		}
//
//		if len(argNameSplit) > 1 {
//			at := argNameSplit[1]
//			if strings.HasSuffix(at, "!") {
//				if ar.defaultVal != "" {
//					log.Fatalf(
//						"Argument '%s' has a default value %s, required args must not have default values",
//						ar.name, ar.defaultVal,
//					)
//				}
//
//				ar.required = true
//			}
//
//			ar.argType = strings.TrimSuffix(
//				strings.TrimSpace(at),
//				"!",
//			)
//		}
//
//		args[ar.name] = ar
//	}
//
//	return args
//}
//
//func (b *Bobfile) convertBody(body string) []vm.Ins {
//	var ins []vm.Ins
//
//	for _, s := range strings.Split(body, "\n") {
//		cleanCmd := cleanLine(s)
//		if cleanCmd == "" || strings.HasPrefix(cleanCmd, "//") {
//			continue
//		}
//
//		switch {
//		case isFn(cleanCmd):
//			ins = append(ins, b.parseBody(cleanCmd)...)
//		case isVar(cleanCmd):
//			varInitIns, err := b.parseVar(cleanCmd)
//			if err != nil {
//				// todo handle err
//				log.Fatal(err)
//			}
//
//			ins = append(
//				ins,
//				varInitIns...,
//			)
//		}
//
//		//ts.cmd = cleanCmd
//		//tasks = append(tasks, ts)
//	}
//
//	return ins
//}
//
//func (b *Bobfile) parseBody(cleanCmd string) []vm.Ins {
//	var ins []vm.Ins
//
//	splits := strings.SplitN(cleanCmd, "(", 2)
//	fnName := splits[0]
//	if len(splits) > 1 {
//		clean := strings.TrimSuffix(splits[1], ")")
//
//		quoteC := 0
//		var args []string
//		var sb strings.Builder
//		for _, c := range clean {
//			stcC := string(c)
//
//			switch {
//			case stcC == `"`:
//				if quoteC == 1 {
//					quoteC = 0
//				} else {
//					quoteC++
//				}
//			case c == ',' && quoteC == 0:
//				args = append(args, strings.TrimSpace(sb.String()))
//				sb.Reset()
//			default:
//				sb.WriteString(stcC)
//			}
//		}
//		// todo refactor kinda duplicate dont like it
//		// also the space is getting removed from the quotes fix
//		args = append(args, strings.TrimSpace(sb.String()))
//
//		for i := len(args) - 1; i >= 0; i-- {
//			ins = append(ins, vm.OVStr(vm.PUSH, args[i]))
//		}
//
//		ins = append(ins, vm.OVInt(vm.PUSH, len(args)))
//	}
//
//	return append(ins, vm.AddFnCall(fnName, false)...)
//}
//
//// function calls start with a valid English alpha and end with ')'
//func isFn(line string) bool {
//	return strings.HasSuffix(line, ")") && isAlpha("")
//}
//
//func (b *Bobfile) parseVersion(line CleanLine) (bool, error) {
//	if !strings.HasPrefix(line.content, VersionPrefix) {
//		return false, nil
//	}
//
//	splits := strings.Split(line.content, VersionPrefix)
//	if len(splits) != 2 {
//		return false, fmt.Errorf("could not parse version line, %s, %s", line.content, splits)
//	}
//
//	s := strings.TrimSpace(splits[1])
//	atoi, err := strconv.Atoi(s)
//	if err != nil {
//		return false, fmt.Errorf("could not convert version to int, %s, %s, %v", line.content, s, err)
//	}
//
//	b.Version = atoi
//	return true, nil
//}
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//// utils
//
//func tokenize(s string) []string {
//	var tokens []string
//	var current strings.Builder
//
//	for _, ch := range s {
//		switch ch {
//		case '+', '-', '*', '/', '(', ')':
//			if current.Len() > 0 {
//				tokens = append(tokens, strings.TrimSpace(current.String()))
//				current.Reset()
//			}
//			tokens = append(tokens, string(ch))
//		case ' ':
//			// skip
//		default:
//			current.WriteRune(ch)
//		}
//	}
//	if current.Len() > 0 {
//		tokens = append(tokens, strings.TrimSpace(current.String()))
//	}
//	return tokens
//}
//
//func toPostfix(s string) []string {
//	var output []string
//	var ops []string
//
//	precedence := map[string]int{
//		"+": 1, "-": 1,
//		"*": 2, "/": 2,
//	}
//
//	tokens := tokenize(s)
//
//	for _, tok := range tokens {
//		switch tok {
//		case "+", "-", "*", "/":
//			for len(ops) > 0 {
//				top := ops[len(ops)-1]
//				if top != "(" && precedence[top] >= precedence[tok] {
//					output = append(output, top)
//					ops = ops[:len(ops)-1]
//				} else {
//					break
//				}
//			}
//			ops = append(ops, tok)
//		case "(":
//			ops = append(ops, tok)
//		case ")":
//			for len(ops) > 0 && ops[len(ops)-1] != "(" {
//				output = append(output, ops[len(ops)-1])
//				ops = ops[:len(ops)-1]
//			}
//			if len(ops) > 0 {
//				ops = ops[:len(ops)-1] // pop the "("
//			}
//		default:
//			output = append(output, tok)
//		}
//	}
//
//	for len(ops) > 0 {
//		output = append(output, ops[len(ops)-1])
//		ops = ops[:len(ops)-1]
//	}
//
//	return output
//}
