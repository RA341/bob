package parser

import (
	"testing"

	"github.com/RA341/bob/vm"
	"github.com/stretchr/testify/require"
)

func TestParser_Expression_Simple(t *testing.T) {
	content := `1 + 2 + (2 + 3)`
	p := runParse(t, content)

	for _, s := range p {
		vv := executeIns(s)
		v := vv.Stack.MustPop()
		require.Equal(t, v.Raw, "8")
	}

	// todo not sure how to handle this
	//content = `(1 + 2)(2 + 3)`
	//runParse(t, content)
}

func TestParser_Expression_Math(t *testing.T) {
	content := `(2 + 3) * 2 + 1`
	p := runParse(t, content)

	var v vm.Value
	for _, s := range p {
		ins := s.Ins()
		vm.NewProgram(ins).Print()

		vv := vm.VM{}
		vv.Start(ins, nil)
		v = vv.Stack.MustPop()
	}

	content = `1 + 2 * (2 + 3)`
	p = runParse(t, content)
	var v2 vm.Value
	for _, s := range p {
		ins := s.Ins()
		vm.NewProgram(ins).Print()

		vv := vm.VM{}
		vv.Start(ins, nil)
		v2 = vv.Stack.MustPop()
	}

	require.Equal(t, v.Raw, "11")
	require.Equal(t, v.Raw, v2.Raw)
}

func TestParser_Var(t *testing.T) {
	content := `var someVar = "test"`
	p1 := runParse(t, content)

	content = `var someVar`
	p := runParse(t, content)
	require.NotEqual(t, p1[0].Ins(), p[0].Ins())

	content = `someVar := "test"`
	p2 := runParse(t, content)

	require.Equal(t, len(p1), 1)
	require.Equal(t, len(p2), len(p1))
	ins := p2[0].Ins()
	i := p1[0].Ins()
	require.Equal(t, i, ins)

	for _, s := range p {
		gf := s.Ins()
		vm.NewProgram(gf).Print()

		vv := vm.VM{}
		vv.Start(gf, nil)
		_ = vv.Stack.MustPop()
	}
}

func runParse(t *testing.T, content string) []Statement {
	lex := RunLexer([]byte(content))
	require.Nil(t, lex.errs)

	p, err := RunParser(lex.tokens)
	require.NoError(t, err)

	return p
}

func executeIns(in Instructions) vm.VM {
	ins := in.Ins()

	vm.NewProgram(ins).Print()

	vv := vm.VM{}
	vv.Start(ins, nil)
	return vv
}
