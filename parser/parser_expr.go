package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/RA341/bob/vm"
)

type Expr interface {
	Str() string
	Ins() []vm.Ins
}

//////////////////////////////////////////////////////////////////////

type ExprNum struct {
	tok Token
}

func (e ExprNum) Ins() []vm.Ins {
	// this token is expected to be a valid num
	atoi, _ := strconv.Atoi(e.tok.Literal)

	return []vm.Ins{
		vm.OVInt(vm.PUSH, atoi),
	}
}

func (e ExprNum) Str() string {
	return e.tok.Lexeme
}

//////////////////////////////////////////////////////////////////////

type ExprString struct {
	tok Token
}

func (e ExprString) Ins() []vm.Ins {
	return []vm.Ins{
		vm.OVStr(vm.PUSH, e.tok.Literal),
	}
}

func (e ExprString) Str() string {
	return e.tok.Lexeme
}

//////////////////////////////////////////////////////////////////////

type ExprGrouping struct {
	expr Expr
}

func (e ExprGrouping) Ins() []vm.Ins {
	// todo sketchy not sure
	return e.expr.Ins()
}

func (e ExprGrouping) Str() string {
	return fmt.Sprintf("(%s)", e.expr.Str())
}

//////////////////////////////////////////////////////////////////////

type ExprBinary struct {
	left     Expr
	operator Token
	right    Expr
}

func (e ExprBinary) Ins() []vm.Ins {
	var op []vm.Ins

	switch e.operator.TokenType {
	case PLUS:
		op = []vm.Ins{
			vm.O(vm.ADD),
		}
	case MINUS:
		op = []vm.Ins{
			vm.O(vm.SUB),
		}
	case STAR:
		op = []vm.Ins{
			vm.O(vm.MUL),
		}
	case SLASH:
		op = []vm.Ins{
			vm.O(vm.DIV),
		}
	case EqualEqual:
		op = []vm.Ins{
			vm.O(vm.EQ),
		}
	case BangEqual:
		op = []vm.Ins{
			vm.O(vm.EQ),
		}
	case Greater:
		op = []vm.Ins{
			vm.O(vm.GT),
		}
	case Less:
		op = []vm.Ins{
			vm.O(vm.LT),
		}
	case LessEqual:
		op = []vm.Ins{
			vm.O(vm.LT),
			vm.O(vm.OR),
			vm.O(vm.EQ),
		}
	case GreaterEqual:
		op = []vm.Ins{
			vm.O(vm.GT),
			vm.O(vm.OR),
			vm.O(vm.EQ),
		}
	default:
		log.Printf("Unsupported binary operator type:\n%s", e.operator.String())
		// todo handle err
		return nil
	}

	lins := e.left.Ins()
	rins := e.right.Ins()
	factors := append(
		lins,
		rins...,
	)

	return append(factors, op...)
}

func (e ExprBinary) Str() string {
	return fmt.Sprintf(
		"%s %s %s",
		e.left.Str(),
		e.operator.Lexeme,
		e.right.Str(),
	)
}

//////////////////////////////////////////////////////////////////////

type ExprLiteral struct {
	tok Token
}

func (e *ExprLiteral) Ins() []vm.Ins {
	return []vm.Ins{
		vm.OVStr(vm.PUSH, e.tok.Literal),
		vm.O(vm.LOAD),
	}
}

func (e *ExprLiteral) Str() string {
	return e.tok.Lexeme
}

//////////////////////////////////////////////////////////////////////

type ExprUnary struct {
	operator Token
	left     Expr
}

func (e *ExprUnary) Ins() []vm.Ins {
	var op vm.Ins

	switch e.operator.TokenType {
	case Bang:
		op = vm.O(vm.NOT)
	//case MINUS: todo
	//	op = vm.O(vm.)
	default:
		log.Printf("Unsupported unary operator type:\n%s", e.operator.String())
		return nil
	}

	return append(e.left.Ins(), op)
}

func (e *ExprUnary) Str() string {
	return fmt.Sprintf(
		"%s%s",
		e.operator.Lexeme,
		e.left.Str(),
	)
}

//////////////////////////////////////////////////////////////////////

type Stmt struct {
}

func (s *Stmt) Ins() []vm.Ins {
	//TODO implement me
	panic("implement me")
}

func (s *Stmt) Str() string {
	//TODO implement me
	panic("implement me")
}

//////////////////////////////////////////////////////////////////////
