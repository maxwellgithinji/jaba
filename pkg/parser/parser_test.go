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

	checkParseError(t, p)

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

func checkParseError(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		t.Errorf("parser error: %s", message)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
       return 1;
	   return 10,
	   return 9992919921
    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParseError(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not have 3 statements, got: %d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got: %T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral() is not'return', got: %s", returnStatement.TokenLiteral())
			continue
		}

	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	checkParseError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement not ast.ExpressionStatement, got: %T", statement)
	}

	identifierExpression, ok := statement.Value.(*ast.Identifier)
	if !ok {
		t.Errorf("expressionStatement.Value not *ast.Identifier, got: %T", statement.Value)
	}

	if identifierExpression.TokenLiteral() != "foobar" {
		t.Errorf("identifierExpression.TokenLiteral() is not 'foobar', got: %s", identifierExpression.TokenLiteral())
	}
}
