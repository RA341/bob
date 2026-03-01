package parser

import (
	"fmt"
)

type Parser struct {
	tokens []Token

	start   int
	current int
}

func NewParser(token []Token) {
	t := Parser{
		tokens: token,
	}

	t.Parse()
}

func (t *Parser) Parse() {
	for !t.isAtEnd() {
		t.start = t.current
		t.ParseToken()
	}
}

func (t *Parser) ParseToken() {
	tok := t.Next()

	fmt.Println(tok.String())
}

func (t *Parser) Peek() Token {
	return t.tokens[t.current]
}

func (t *Parser) Match(tt ...TokenType) bool {
	for _, ty := range tt {
		if t.Check(ty) {
			t.Next()
			return true
		}
	}

	return false
}

func (t *Parser) Check(tt TokenType) bool {
	if t.isAtEnd() {
		return false
	}
	return t.Peek().TokenType == tt
}

func (t *Parser) Prev() Token {
	return t.tokens[t.current-1]
}

func (t *Parser) Next() Token {
	s := t.Peek()
	t.current += 1
	return s
}

func (t *Parser) isAtEnd() bool {
	return t.Peek().TokenType == EOF
}
