package parser

import (
	"testing"

	"github.com/RA341/bob/vm"
	"github.com/stretchr/testify/require"
)

func TestParser_Expression(t *testing.T) {
	content := `1 + 2 + (2 + 3)`
	runParse(t, content)

	// todo not sure how to handle this
	//content = `(1 + 2)(2 + 3)`
	//runParse(t, content)
}

func TestParser_Expression_vm(t *testing.T) {
	content := `1 + 2 * (2 + 3)`
	p := runParse(t, content)

	ins := p.Ins()
	vv := vm.VM{}
	vv.Start(ins, nil)

	v := vv.Stack.MustPop()
	require.Equal(t, v.Raw, "8")
}

func runParse(t *testing.T, content string) Expr {
	lex := RunLexer([]byte(content))
	require.Nil(t, lex.errs)

	p, err := RunParser(lex.tokens)
	require.NoError(t, err)

	t.Log(p.Str())

	return p
}
