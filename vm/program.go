package vm

import (
	"fmt"
	"strings"
)

type Program struct {
	input []Ins
}

func NewProgram(in []Ins) *Program {
	return &Program{
		input: in,
	}
}

func (p *Program) Get() []Ins {
	return p.input
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, in := range p.input {
		sb.WriteString(in.String() + "\n")
	}

	return sb.String()
}

func (p *Program) Print() {
	for _, in := range p.input {
		fmt.Println(in)
	}
}

func (p *Program) Add(ins ...Ins) {
	p.input = append(p.input, ins...)
}

// AddExpr assigns a variable a value
func AddExpr(varName string, expr ...Ins) []Ins {
	return append(
		expr,
		[]Ins{
			OVStr(PUSH, varName),
			O(STORE),
		}...,
	)
}

// AddVar instructions to register a var
func AddVar(name string, value Value) []Ins {
	return []Ins{
		OV(PUSH, value),
		OVStr(PUSH, name),
		O(STORE),
	}
}

// LoadVar generates instructions to get a var
func LoadVar(name string) []Ins {
	return []Ins{
		OVStr(LOAD, name),
	}
}

func AddFnCall(name string, variadic bool, args ...Value) []Ins {
	var argIns []Ins
	for _, arg := range args {
		argIns = append(
			argIns,
			OV(PUSH, arg),
		)
	}

	var fnCall = []Ins{
		OVStr(PUSH, name),
		O(CALL),
	}

	if variadic {
		fnCall = append(
			[]Ins{
				OVInt(PUSH, len(args)),
			},
			fnCall...,
		)
	}

	return append(
		argIns,
		fnCall...,
	)
}
