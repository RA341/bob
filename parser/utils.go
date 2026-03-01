package parser

import (
	"log"
	"strings"
)

const VarPrefix = "@"

func isVar(line string) bool {
	return strings.HasPrefix(line, VarPrefix)
}

// cleanLine removes new lines and spaces at the start of a line
func cleanLine(spl string) string {
	return strings.Trim(strings.TrimSpace(spl), "\n")
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
