package parser

import (
	"github.com/RA341/bob/vm"
)

type Instructions interface {
	Ins() []vm.Ins
}

type Statement interface {
	Instructions
}

//////////////////////////////////////////////////////////////////////

type StmtExpr struct {
	exp Expr
}

func (e StmtExpr) Ins() []vm.Ins {
	return e.exp.Ins()
}

//////////////////////////////////////////////////////////////////////

type StmtPrint struct {
	expr Expr
}

func (p StmtPrint) Ins() []vm.Ins {
	return append(
		p.expr.Ins(),
		vm.OVStr(vm.PUSH, "print"),
		vm.O(vm.CALL),
	)
}

//////////////////////////////////////////////////////////////////////

type StmtVar struct {
	Initializer Expr
	Identifier  Token
}

func (s StmtVar) Ins() []vm.Ins {
	initializer := []vm.Ins{vm.OVNil(vm.PUSH)}
	if s.Initializer != nil {
		initializer = s.Initializer.Ins()
	}

	return append(
		initializer,
		vm.OVStr(vm.PUSH, s.Identifier.Literal),
		vm.O(vm.STORE),
	)
}
