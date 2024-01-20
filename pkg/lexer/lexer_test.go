package lexer

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Assign, "="},
		{token.Plus, "+"},
		{token.LParen, "("},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Comma, ","},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range test {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected = %s, got %s", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token literal. expected = %s, got %s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
