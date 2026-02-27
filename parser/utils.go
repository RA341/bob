package parser

import "strings"

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
