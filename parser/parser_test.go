package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	content := `
cond1 = value
cond2 = value2

if cond1 == cond2 {
	print("cond1 is equal")
}
`

	lex := NewLexer([]byte(content))
	require.Nil(t, lex.errs)

	NewParser(lex.tokens)
}
