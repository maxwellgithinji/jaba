package lexer

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/token"
)

func TestNextTokenOperatorsDelimitersKeywords(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected = %s, got %s", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token literal. expected = %s, got %s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenJabaProgram(t *testing.T) {
	input := `let foo = 1;
    let bar = 3;

    let add = fn (foo, bar) {
        foo + bar;
    };

    let result = add(foo, bar);
    `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, "="},
		{token.INTEGER, "1"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "bar"},
		{token.ASSIGN, "="},
		{token.INTEGER, "3"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "foo"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "bar"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "bar"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "foo"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "bar"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected = %s, got %s", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong token literal. expected = %s, got %s", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
