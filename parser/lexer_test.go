package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Lex_Basic(t *testing.T) {
	// basic operators
	content := `({}) < <= >= > + - * , : ! != ==`
	lex := RunLexer([]byte(content))
	expected := []TokenType{
		LPAREN, LCURLY, RCURLY, RPAREN,
		Less, LessEqual, GreaterEqual, Greater,
		PLUS, MINUS, STAR, COMMA, Colon,

		Bang, BangEqual, EqualEqual,

		EOF,
	}
	require.Nil(t, lex.errs)
	for i, r := range lex.tokens {
		require.Equal(t, expected[i], r.TokenType)
	}

	// comments
	content = `//({}) < <= >= > + - * , : ! != ==`
	lex = RunLexer([]byte(content))
	expected = []TokenType{EOF}
	require.Nil(t, lex.errs)
	require.Equal(t, len(expected), len(lex.tokens))

	// strings
	content = `"//({}) < <= >= > + - * , : ! != =="`
	lex = RunLexer([]byte(content))
	expected = []TokenType{String, EOF}
	require.Nil(t, lex.errs)
	require.Equal(t, len(expected), len(lex.tokens))
	require.Equal(t, expected[0].String(), lex.tokens[0].TokenType.String())
	require.Equal(t, content[1:len(content)-1], lex.tokens[0].Literal)

	content = `if for else or and`
	lex = RunLexer([]byte(content))
	expected = []TokenType{If, For, Else, Or, And, EOF}
	require.Nil(t, lex.errs)
	for i, r := range lex.tokens {
		require.Equal(t, expected[i], r.TokenType)
	}
}

func TestLexContents(t *testing.T) {
	content := `@as=aasdasd
@sd=1

hello2(
	user!, 
	other2:,
) {
    workdir core
    print("hello ${user}")
    print(as)
}`

	lex := RunLexer([]byte(content))
	require.Nil(t, lex.errs)

	for _, r := range lex.tokens {
		fmt.Println(r.String())
	}

	t.Log(lex.tokens)
}
