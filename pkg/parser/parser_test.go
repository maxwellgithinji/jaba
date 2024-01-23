package parser

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
		let x = 1;
		let y = 12;
		let foo = 123;
		let bar = 1;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 4 {
		for i := 0; i < len(program.Statements); i++ {
			t.Logf("Statement %d: %v", i, program.Statements[i])
		}
		t.Fatalf("ParseProgram() returned %d statements instead of 4", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
		{"bar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatements(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral is not 'let', got %q ", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s is not *ast.LetStatement, got %T", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Fatalf("letStatement.Name.Value is not %s, got %s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Fatalf("letStatement.Name.TokenLiteral() is not %s, got %s", name, letStatement.Name.TokenLiteral())
		return false
	}
	return true
}
