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

	lex := RunLexer([]byte(content))
	require.Nil(t, lex.errs)

	RunParser(lex.tokens)
}
