package parser

import (
	"fmt"
	"log"
)

//go:generate go run github.com/dmarkham/enumer@latest -type=TokenType -output=gen_enum_token.go
type TokenType int

const (
	LCURLY TokenType = iota
	RCURLY
	LPAREN
	RPAREN

	COMMA
	COLON

	Equal
	EqualEqual
	ColonEqual

	Greater
	GreaterEqual
	Less
	LessEqual

	BangEqual
	Bang

	SLASH
	PLUS
	MINUS
	STAR

	Or
	And
	If
	Else
	For

	AT
	Identifier
	Literal

	Num
	String

	EOF
)

type Token struct {
	TokenType TokenType
	Literal   string
	Lexeme    string
	Line      int
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"L%d: T=%s LIT=%q LX=%q",
		t.Line,
		t.TokenType.String(),
		t.Literal,
		t.Lexeme,
	)
}

type Lexer struct {
	tokens []Token

	errs     []error
	contents []byte

	start   int
	current int
	line    int
}

func RunLexer(content []byte) *Lexer {
	lex := Lexer{
		contents: content,
		line:     1,
	}

	lex.Parse()
	return &lex
}

func (l *Lexer) Parse() {
	for !l.isAtEnd() {
		// We are at the beginning of the next lexeme.
		l.start = l.current
		l.ScanToken()
	}

	l.AddToken(EOF)
}

func (l *Lexer) ScanToken() {
	ch := l.Next()
	if ch == "" {
		log.Println("Empty char, THIS SHOULD NEVER HAPPEN")
		return
	}
	switch ch {
	case "{":
		l.AddToken(LCURLY)
	case "}":
		l.AddToken(RCURLY)
	case "(":
		l.AddToken(LPAREN)
	case ")":
		l.AddToken(RPAREN)
	case "+":
		l.AddToken(PLUS)
	case "-":
		l.AddToken(MINUS)
	case "*":
		l.AddToken(STAR)
	case ",":
		l.AddToken(COMMA)
	case ":":
		l.addIfHasNext("=", ColonEqual, COLON)
	case "@":
		l.AddToken(AT)
	case "!":
		l.addIfHasNext("=", BangEqual, Bang)
	case "=":
		l.addIfHasNext("=", EqualEqual, Equal)
	case "<":
		l.addIfHasNext("=", LessEqual, Less)
	case ">":
		l.addIfHasNext("=", GreaterEqual, Greater)
	case "/":
		if l.HasNext("/") {
			// comments till the EOL
			for !l.isAtEnd() && l.Peek() != "\n" {
				l.Next()
			}
		} else {
			l.AddToken(SLASH)
		}
	case `"`:
		l.String()
	case " ", "\t", "\r":
	case "\n":
		l.line++
	default:
		if isNum(ch) {
			l.number(ch)
			break
		}

		if isAlpha(ch) {
			l.identifier()
			break
		}
		l.errs = append(l.errs, fmt.Errorf(
			"unexpected char: %s at line %d, Row:%d",
			ch,
			l.line,
			l.current,
		))
		break
	}

}

func (l *Lexer) addIfHasNext(next string, exists TokenType, notExists TokenType) {
	if l.HasNext(next) {
		l.AddToken(exists)
	} else {
		l.AddToken(notExists)
	}
}

func (l *Lexer) Next() string {
	s := string(l.contents[l.current])
	l.current += 1
	return s
}

func (l *Lexer) Peek() string {
	return l.PeekN(0)
}

func (l *Lexer) String() {
	for l.Peek() != `"` && !l.isAtEnd() {
		if l.Peek() == "\n" {
			l.line++
		}
		l.Next()
	}

	if l.isAtEnd() {
		l.errs = append(l.errs, fmt.Errorf("unterminated string at line %d", l.line))
		return
	}

	// closing "
	l.Next()

	// consume without trailing and starting "
	strLit := l.contents[l.start+1 : l.current-1]
	l.AddTokenLit(String, string(strLit))
}

func (l *Lexer) AddToken(tt TokenType) {
	l.AddTokenLit(tt, "")
}

func (l *Lexer) AddTokenLit(tt TokenType, literal string) {
	tokStr := l.contents[l.start:l.current]

	tok := Token{
		TokenType: tt,
		Literal:   literal,
		Lexeme:    string(tokStr),
		Line:      l.line,
	}

	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) PeekN(n int) string {
	peekIdx := l.current + n
	if l.isAtEndN(peekIdx) {
		return ""
	}
	return string(l.contents[peekIdx])
}

func (l *Lexer) HasNext(s string) bool {
	next := l.PeekN(0)
	if next == "" || next != s {
		return false
	}

	l.current++
	return true
}

func (l *Lexer) isAtEnd() bool {
	return l.isAtEndN(l.current)
}

func (l *Lexer) isAtEndN(n int) bool {
	return n >= len(l.contents)
}

func (l *Lexer) number(ch string) {
	l.AddTokenLit(Num, ch)
}

var keywords = map[string]TokenType{
	"and":  And,
	"or":   Or,
	"for":  For,
	"if":   If,
	"else": Else,
}

func (l *Lexer) identifier() {
	for isAlphaNum(l.Peek()) {
		l.Next()
	}

	iden := l.contents[l.start:l.current]
	keyword, ok := keywords[string(iden)]
	if !ok {
		keyword = Identifier
	}

	l.AddToken(keyword)
}

func isAlphaNum(r string) bool {
	return isAlpha(r) || isNum(r)
}

func isNum(r string) bool {
	return r >= "0" && r <= "9"
}

func isAlpha(r string) bool {
	return (r >= "a" && r <= "z") ||
		(r >= "A" && r <= "Z") ||
		r == "_"
}
