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

func RunParser(token []Token) (Expr, error) {
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

func (p *Parser) Parse() Expr {
	return p.expression()
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
			operator: p.Prev(),
			right:    p.comparison(),
		}
	}

	return left
}

// comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	term1 := p.term()

	for p.Match(Less, LessEqual, Greater, GreaterEqual) {
		op := p.Prev()
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
		op := p.Prev()
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
		un2 := p.unary()
		op := p.Prev()
		un = ExprBinary{
			left:     un,
			operator: op,
			right:    un2,
		}
	}

	return un
}

// unary -> ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	if p.Match(Bang, MINUS) {
		op := p.Prev()
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
			tok: p.Prev(),
		}
	}

	if p.Match(String) {
		return ExprString{
			tok: p.Prev(),
		}
	}

	if p.Match(Literal) {
		return &ExprLiteral{
			tok: p.Prev(),
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

func (p *Parser) declaration() (Stmt, error) {
	if p.Match(Identifier, ColonEqual) {
		return p.varDec(Identifier)
	}

	//return t.expression()
	return Stmt{}, nil
}

func (p *Parser) sync() {

}

func (p *Parser) varDec(identifier TokenType) (Stmt, error) {
	return Stmt{}, nil
}

func (p *Parser) consume(tok TokenType, message string) {
	peek := p.Peek()
	if peek.TokenType == tok {
		p.Next()
		return
	}

	p.errs = append(
		p.errs,
		NewParseErr(peek, message),
	)
}

func (p *Parser) Peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) MatchExtract(tt ...TokenType) (Token, bool) {
	if p.Match(tt...) {
		return p.Prev(), true
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

func (p *Parser) Prev() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) Next() Token {
	s := p.Peek()
	p.current += 1
	return s
}

func (p *Parser) isAtEnd() bool {
	return p.Peek().TokenType == EOF
}
