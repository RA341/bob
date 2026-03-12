package parser

import (
	"errors"
	"fmt"

	"github.com/RA341/bob/vm"
)

type ParseErr struct {
	tok Token
	msg string
}

func NewParseErr(tok Token, msg string) ParseErr {
	return ParseErr{
		tok: tok,
		msg: msg,
	}
}

func (e *ParseErr) Err() error {
	if e.tok.TokenType == EOF {
		return fmt.Errorf("%d at end %s", e.tok.Line, e.msg)
	}

	return fmt.Errorf("%d at '%s' %s", e.tok.Line, e.tok.Lexeme, e.msg)
}

type Parser struct {
	tokens []Token
	errs   []ParseErr

	current      int
	instructions []vm.Ins
}

func RunParser(token []Token) ([]Statement, error) {
	p := NewParser(token)
	result := p.Parse()

	return result, p.Err()
}

func NewParser(token []Token) *Parser {
	t := Parser{
		tokens: token,
	}

	return &t
}

func (p *Parser) Err() error {
	var err error
	for _, er := range p.errs {
		err = errors.Join(err, er.Err())
	}

	return err
}

func (p *Parser) Parse() []Statement {
	var statements []Statement

	for !p.isAtEnd() {
		statements = append(
			statements,
			p.deceleration(),
		)
	}

	return statements
}

func (p *Parser) deceleration() Statement {
	if p.Match(Var) {
		return p.varDec()
	}

	if p.MatchOffset(1, ColonEqual) {
		return p.varImplicitDec()
	}

	return p.statement()
}

// var using 'var someVar = dec'
func (p *Parser) varDec() Statement {
	iden := p.consume(Identifier, "Expected identifier after 'var")

	var init Expr
	if p.Match(Equal) {
		init = p.expression()
	}

	return StmtVar{
		Initializer: init,
		Identifier:  iden,
	}
}

// var using the ':='
func (p *Parser) varImplicitDec() Statement {
	identifier := p.prevN(2)
	if identifier.TokenType != Identifier {
		p.addErr(ParseErr{
			tok: identifier,
			msg: "expected identifier before ':='",
		})

		return nil
	}

	init := p.expression()

	return StmtVar{
		Identifier:  identifier,
		Initializer: init,
	}
}

func (p *Parser) statement() Statement {
	if p.Match(Print) {
		return p.printStmt()
	}

	return p.exprStmt()
}

func (p *Parser) printStmt() Statement {
	expr := p.expression()
	return StmtPrint{
		expr: expr,
	}
}

func (p *Parser) exprStmt() Statement {
	expr := p.expression()
	return StmtExpr{
		exp: expr,
	}
}

// expression -> equality ;
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	left := p.comparison()

	for p.Match(EqualEqual, BangEqual) {
		left = ExprBinary{
			left:     left,
			operator: p.prev(),
			right:    p.comparison(),
		}
	}

	return left
}

// comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	term1 := p.term()

	for p.Match(Less, LessEqual, Greater, GreaterEqual) {
		op := p.prev()
		term2 := p.term()

		term1 = ExprBinary{
			left:     term1,
			operator: op,
			right:    term2,
		}
	}

	return term1
}

// term  -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	fac := p.factor()

	for p.Match(PLUS, MINUS) {
		op := p.prev()
		fac2 := p.factor()

		fac = ExprBinary{
			left:     fac,
			operator: op,
			right:    fac2,
		}
	}

	return fac
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	un := p.unary()

	for p.Match(SLASH, STAR) {
		// always get the op first before getting right
		op := p.prev()
		un = ExprBinary{
			left:     un,
			operator: op,
			right:    p.unary(),
		}
	}

	return un
}

// unary -> ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	if p.Match(Bang, MINUS) {
		op := p.prev()
		return &ExprUnary{
			operator: op,
			left:     p.unary(),
		}
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
// true/false is unsupported for now
func (p *Parser) primary() Expr {
	if p.Match(Num) {
		return ExprNum{
			tok: p.prev(),
		}
	}

	if p.Match(String) {
		return ExprString{
			tok: p.prev(),
		}
	}

	if p.Match(Literal) {
		return &ExprLiteral{
			tok: p.prev(),
		}
	}

	if p.Match(Identifier) {
		return ExprVar{
			tok: p.prev(),
		}
	}

	if p.Match(LPAREN) {
		group := p.expression()
		p.consume(
			RPAREN,
			"Expected closing ')' after expression",
		)

		return ExprGrouping{
			expr: group,
		}
	}

	p.errs = append(
		p.errs,
		NewParseErr(p.Peek(), "Expected expression"),
	)

	// todo this is sketchy not sure how to handle it
	return nil
}

func (p *Parser) sync() {

}

func (p *Parser) consume(tok TokenType, message string) Token {
	peek := p.Peek()
	if peek.TokenType == tok {
		return p.Next()
	}

	p.errs = append(
		p.errs,
		NewParseErr(peek, message),
	)

	return Token{}
}

func (p *Parser) Peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) MatchOffset(offset int, match ...TokenType) bool {
	for _, ty := range match {
		if p.CheckN(ty, offset) {
			// include the current the token
			for range offset + 1 {
				p.Next()
			}
			return true
		}
	}

	return false
}

func (p *Parser) CheckN(match TokenType, offset int) bool {
	of := p.current + offset
	if of >= len(p.tokens) {
		return false
	}

	tok := p.tokens[of]
	if tok.TokenType == EOF {
		return false
	}

	return tok.TokenType == match
}

func (p *Parser) MatchExtract(tt ...TokenType) (Token, bool) {
	if p.Match(tt...) {
		return p.prev(), true
	}

	return Token{}, false
}

func (p *Parser) Match(tt ...TokenType) bool {
	for _, ty := range tt {
		if p.Check(ty) {
			p.Next()
			return true
		}
	}

	return false
}

func (p *Parser) Check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.Peek().TokenType == tt
}

func (p *Parser) prev() Token {
	return p.prevN(1)
}

func (p *Parser) prevN(off int) Token {
	return p.tokens[p.current-off]
}

func (p *Parser) Next() Token {
	s := p.Peek()
	p.current += 1
	return s
}

func (p *Parser) isAtEnd() bool {
	return p.Peek().TokenType == EOF
}

func (p *Parser) addErr(pe ParseErr) {
	p.errs = append(p.errs, pe)
}
