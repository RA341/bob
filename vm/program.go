package vm

import "fmt"

type Program struct {
	input []Ins
}

func (p *Program) Get() []Ins {
	return p.input
}

func (p *Program) Print() {
	for _, in := range p.input {
		fmt.Println(in)
	}

}

func (p *Program) AddGlobalVar(name string, value Value) {
	varLoad := []Ins{
		OV(PUSH, value),
		OVStr(PUSH, name),
		O(STORE),
	}

	p.input = append(p.input, varLoad...)
}

// GetVarIns generates instructions to get a var
func (p *Program) GetVarIns(name string) []Ins {
	return []Ins{
		OVStr(LOAD, name),
	}
}

func (p *Program) AddGlobalExpr(name string, value ...Ins) {
	varLoad := append(
		value,
		[]Ins{
			OVStr(PUSH, name),
			O(STORE),
		}...,
	)

	p.input = append(p.input, varLoad...)
}
