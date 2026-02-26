package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func ParseFromFile(bobFile *Bobfile, bobFilePath string) {
	fBytes, err := os.ReadFile(bobFilePath)
	if err != nil {
		log.Fatal("Could not read file", err)
	}

	ParseFromBytes(bobFile, fBytes)
}

func ParseFromBytes(bobFile *Bobfile, fBytes []byte) {
	contents := string(fBytes)
	contents = strings.TrimSpace(contents)

	lines := splitAndStrip(contents)

	bobFile.Vars = make(map[string]string)
	bobFile.Cmds = make(map[string]Cmd)

	//fmt.Println("Contents:")
	for _, line := range lines {
		//fmt.Printf("line %d: %s\n", i, line.content)

		parseLine(bobFile, line)
	}

}

const VersionPrefix = "version:"

func parseLine(b *Bobfile, line CleanLine) {
	if parseVersion(b, line) {
		return
	}

	if parseGlobalVar(b, line) {
		return
	}

	if parseCmd(b, line) {
		return
	}
}

func isInsideBraces(s, target string) bool {
	depth := 0
	idx := strings.Index(s, target)
	if idx == -1 {
		return false
	}
	for _, ch := range s[:idx] {
		switch ch {
		case '{':
			depth++
		case '}':
			depth--
		}
	}
	return depth > 0
}

func parseGlobalVar(b *Bobfile, line CleanLine) bool {
	if isInsideBraces(line.content, VarPrefix) {
		// this is a local var
		return false
	}

	key, value := parseVar(line.content)
	if key == "" || value == "" {
		return false
	}

	b.Vars[key] = value
	return true
}

const VarPrefix = "@"

func isVar(line string) bool {
	return strings.HasPrefix(line, VarPrefix)
}

func parseVar(line string) (key string, val string) {
	if !isVar(line) {
		return "", ""
	}

	line = strings.TrimPrefix(line, VarPrefix)
	split := strings.Split(line, "=")
	if len(split) < 2 {
		return "", ""
	}

	key = strings.TrimSpace(split[0])
	val = strings.TrimSpace(split[1])

	return key, val
}

func parseCmd(b *Bobfile, line CleanLine) bool {
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
	UNUSED(retVal)

	returnTypeSet := false

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

	cmd.args = convertTokToArgs(args)
	cmd.tasks, cmd.vars = convertBodyToTasks(body)
	cmd.name = cmdName
	b.Cmds[cmdName] = cmd

	return true
}

func convertBodyToTasks(body string) ([]Task, VarMap) {
	var tasks []Task
	var varMap = VarMap{}

	for _, s := range strings.Split(body, "\n") {
		cleanCmd := cleanLine(s)
		if cleanCmd == "" || strings.HasPrefix(cleanCmd, "//") {
			continue
		}

		var ts Task
		if isVar(cleanCmd) {
			key, val := parseVar(cleanCmd)
			if key == "" || val == "" {
				// todo handle this
				continue
			}

			varMap.Add(key, val)

			continue
		}

		ts.cmd = cleanCmd
		tasks = append(tasks, ts)
	}

	return tasks, varMap
}

// cleanLine removes new lines and spaces at the start of a line
func cleanLine(spl string) string {
	return strings.Trim(strings.TrimSpace(spl), "\n")
}

func convertTokToArgs(rawArgs string) map[string]Arg {
	args := make(map[string]Arg)
	if rawArgs == "" {
		return args
	}

	// in: user: str!, \notherP: sd,
	split := strings.Split(rawArgs, ",")
	for _, arg := range split {
		if arg == "" {
			continue
		}

		var ar Arg

		arg = cleanLine(arg)
		defV := strings.Split(arg, "=")

		// has default
		if len(defV) > 1 {
			ar.defaultVal = cleanLine(defV[1])
		}

		segs := strings.Split(defV[0], ":")

		// does not have type info
		// [user]
		if len(segs) > 0 {
			argName := segs[0]
			if _, ok := args[argName]; ok {
				log.Fatalf("Argument '%s' already defined", argName)
			}

			ar.name = argName
		}

		if len(segs) > 1 {
			at := segs[1]
			if strings.HasSuffix(at, "!") {
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

func parseVersion(b *Bobfile, line CleanLine) bool {
	if !strings.HasPrefix(line.content, VersionPrefix) {
		return false
	}

	splits := strings.Split(line.content, VersionPrefix)
	if len(splits) != 2 {
		log.Fatal("Could not parse version line", line, splits)
	}
	s := strings.TrimSpace(splits[1])
	atoi, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Could not convert version to int", line, s, err)
	}

	b.Version = atoi
	return true
}

type CleanLine struct {
	lineNum int
	content string
}

func splitAndStrip(contents string) []CleanLine {
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
