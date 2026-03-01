package expr

type Expr interface {
	Visit(expr Expr)
}
