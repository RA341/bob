package parser

type Expr interface {
	Visit(expr Expr)
}

type ExprBool struct {
	left     Expr
	operator Token
	right    Expr
}

func (e ExprBool) Visit(expr Expr) {}

type ExprLiteral struct {
	tok Token
}

func (e ExprLiteral) Visit(expr Expr) {}

type ExprUnary struct {
	operator Token
	left     Expr
}

func (e ExprUnary) Visit(expr Expr) {}

type Grouping struct {
	exp Expr
}

func (e ExprUnary) Grouping(expr Expr) {}
