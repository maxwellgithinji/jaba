package parser

import (
	"fmt"
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/lexer"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 1;", "x", 1},
		{"let y = 12;", "y", 12},
		{"let foo = 123;", "foo", 123},
		{"let bar = 1;", "bar", 1},
		{"let z = true;", "z", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatements(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}

	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 1;", 1},
		{"return 10;", 10},
		{"return 9992919921;", 9992919921},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got: %T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral() is not'return', got: %s", returnStatement.TokenLiteral())
			continue
		}

		if testLiteralExpression(t, returnStatement.Value, tt.expectedValue) {
			return
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

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

	literal, ok := statement.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expressionStatement.Value not *ast.IntegerLiteral, got: %T", statement.Value)
	}

	if literal.Value != 5 {
		t.Errorf("literal.TokenLiteral() is not 5, got: %s", literal.TokenLiteral())
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not %s, got: %s", "5", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		checkParseError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
		}

		expression, ok := statement.Value.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("statement.Value is not ast.PrefixExpression, got: %T", statement.Value)
		}

		if expression.Operator != tt.operator {
			t.Errorf("expression.Operator is not %s, got: %s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5 ", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false != false", false, "!=", false},
		{"false != true", false, "!=", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		checkParseError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
		}

		if !testInfixExpression(t, statement.Value, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b -c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a +  b * c  + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 3][b * c] * d",
			"((a * ([1, 2, 3, 3][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		checkParseError(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, actual)
		}

	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
		}

		boolean, ok := statement.Value.(*ast.Boolean)
		if !ok {
			t.Errorf("statement.Value is not ast.Boolean, got: %T", statement.Value)
		}
		if boolean.Value != tt.expected {
			t.Errorf("boolean.Value is not %t, got: %t", tt.expected, boolean.Value)
		}
	}
}

// Test Helper function
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

func TestIfExpression(t *testing.T) {
	input := "if (x < y ){ x }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
	}

	expression, ok := statement.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Value is not ast.IfExpression, got: %T", statement.Value)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("expression.Consequence.Statements expected 1 statements, got: %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not ast.ExpressionStatement, got: %T", consequence)
	}

	if !testIdentifier(t, consequence.Value, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements was not nil, got: %+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else {y}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements expected 1 statements, got: %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
	}

	expression, ok := statement.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Value is not ast.IfExpression, got: %T", statement.Value)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("expression.Consequence.Statements expected 1 statements, got: %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not ast.ExpressionStatement, got: %T", consequence)
	}

	if !testIdentifier(t, consequence.Value, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("len expression.Alternative.Statements was not 1, got: %+v", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Alternative.Statements[0] is not ast.ExpressionStatement, got: %T", alternative)
	}

	if !testIdentifier(t, alternative.Value, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements expected 1 statement, got: %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
	}

	function, ok := statement.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.FunctionLiteral, got: %T", statement.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters expected 2 parameters, got: %d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements expected 1 statement, got: %d", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement, got: %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Value, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {}", []string{}},
		{"fn(x) {}", []string{"x"}},
		{"fn(x, y, z) {}", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		P := New(l)

		program := P.ParseProgram()

		checkParseError(t, P)

		statement := program.Statements[0].(*ast.ExpressionStatement)

		function := statement.Value.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf("function.Parameters expected %d parameters, got: %d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, identifier := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}

	}
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression, got %T", expression)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("operatorExpression.Operator is not %s, got %s", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("type of expression not supported, got %T", expression)

	return false
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	fmt.Println(expression, value)
	integer, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression is not an ast.IntegerLiteral, got %T", expression)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() is not %d, got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, vale bool) bool {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression is not an ast.Boolean, got %T", expression)
		return false
	}

	if boolean.Value != vale {
		t.Errorf("boolean.Value is not %t, got %t", vale, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", vale) {
		t.Errorf("boolean.TokenLiteral() is not %t, got %s", vale, boolean.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not ast.Identifier, got %T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s, got %s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() is not %s, got %s", value, identifier.TokenLiteral())
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

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements expected 1 statement, got: %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got: %T", statement)
	}

	callExpression, ok := statement.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement.Value is not ast.CallExpression, got: %T", statement.Value)
	}

	if !testIdentifier(t, callExpression.Function, "add") {
		return
	}

	if len(callExpression.Arguments) != 3 {
		t.Fatalf("callExpression.Arguments expected 3 arguments, got: %d", len(callExpression.Arguments))
	}

	testLiteralExpression(t, callExpression.Arguments[0], 1)
	testInfixExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpression.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedArguments  []string
	}{
		{
			input:              "add();",
			expectedIdentifier: "add",
			expectedArguments:  []string{},
		},
		{
			input:              "add(1);",
			expectedIdentifier: "add",
			expectedArguments:  []string{"1"},
		},

		{
			input:              "add(1, 2 * 3, 4 + 5);",
			expectedIdentifier: "add",
			expectedArguments:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		P := New(l)

		program := P.ParseProgram()

		checkParseError(t, P)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		callExpression, ok := statement.Value.(*ast.CallExpression)
		if !ok {
			t.Fatalf("statement.Value is not ast.CallExpression, got: %T", statement.Value)
		}

		if !testIdentifier(t, callExpression.Function, tt.expectedIdentifier) {
			return
		}

		if len(callExpression.Arguments) != len(tt.expectedArguments) {
			t.Fatalf("callExpression.Arguments expected %d arguments, got: %d", len(tt.expectedArguments), len(callExpression.Arguments))
		}

		for i, argument := range tt.expectedArguments {
			if callExpression.Arguments[i].String() != argument {
				t.Errorf("argument %d is  not %s, got %s", i, argument, callExpression.Arguments[i].String())
			}
		}

	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	stringLiteral, ok := statement.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.StringLiteral, got: %T", statement.Value)
	}

	if stringLiteral.Value != "hello world" {
		t.Errorf("stringLiteral.Value is not %s, got %s", "hello world", stringLiteral.Value)
		return
	}

}

func TestParsingArrayLiteral(t *testing.T) {
	input := "[1 , 2 * 2, 3 + 3]"

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	arrayLiteral, ok := statement.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.ArrayLiteral, got: %T", statement.Value)
	}

	if len(arrayLiteral.Elements) != 3 {
		t.Fatalf("arrayLiteral.Elements expected 3 elements, got: %d", len(arrayLiteral.Elements))
	}

	testIntegerLiteral(t, arrayLiteral.Elements[0], 1)
	testInfixExpression(t, arrayLiteral.Elements[1], 2, "*", 2)
	testInfixExpression(t, arrayLiteral.Elements[2], 3, "+", 3)

}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1 + 1]`

	l := lexer.New(input)
	P := New(l)
	program := P.ParseProgram()
	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	indexExpression, ok := statement.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("statement.Value is not ast.IndexExpression, got: %T", statement.Value)
	}

	if !testIdentifier(t, indexExpression.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExpression.Index, 1, "+", 1) {
		return
	}

}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	hashLiteral, ok := statement.Value.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.HashLiteral, got: %T", statement.Value)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Fatalf("hashLiteral.Pairs expected 3 pairs, got: %d", len(hashLiteral.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("key is not ast.StringLiteral, got: %T", key)
		}

		expectedValues := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValues)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	hashLiteral, ok := statement.Value.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.HashLiteral, got: %T", statement.Value)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hashLiteral.Pairs) != len(expected) {
		t.Fatalf("hashLiteral.Pairs expected %d pairs, got: %d", len(expected), len(hashLiteral.Pairs))
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("key is not ast.IntegerLiteral, got: %T", key)
		}

		expectedValues := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValues)
	}

}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	hashLiteral, ok := statement.Value.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.HashLiteral, got: %T", statement.Value)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hashLiteral.Pairs) != len(expected) {
		t.Fatalf("hashLiteral.Pairs expected 2 pairs, got: %d", len(hashLiteral.Pairs))
	}

	for key, value := range hashLiteral.Pairs {
		boolean, ok := key.(*ast.Boolean)
		if !ok {
			t.Fatalf("key is not ast.Boolean, got: %T", key)
		}

		expectedValues := expected[boolean.String()]

		testIntegerLiteral(t, value, expectedValues)
	}

}

func TestParsingEmptyHashLiterals(t *testing.T) {
	input := "{}"

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	hashLiteral, ok := statement.Value.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.HashLiteral, got: %T", statement.Value)
	}

	if len(hashLiteral.Pairs) != 0 {
		t.Fatalf("hashLiteral.Pairs expected 0 pairs, got: %d", len(hashLiteral.Pairs))
	}

}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one":  0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)

	P := New(l)

	program := P.ParseProgram()

	checkParseError(t, P)

	statement := program.Statements[0].(*ast.ExpressionStatement)

	hashLiteral, ok := statement.Value.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement.Value is not ast.HashLiteral, got: %T", statement.Value)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Fatalf("hashLiteral.Pairs expected 3 pairs, got: %d", len(hashLiteral.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one":   func(e ast.Expression) { testInfixExpression(t, e, 0, "+", 1) },
		"two":   func(e ast.Expression) { testInfixExpression(t, e, 10, "-", 8) },
		"three": func(e ast.Expression) { testInfixExpression(t, e, 15, "/", 5) },
	}

	for key, value := range hashLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("key is not ast.StringLiteral, got: %T", key)
			continue
		}

		test, ok := tests[literal.String()]
		if !ok {
			t.Fatalf("key is not a valid test, got: %s", literal.String())
			continue
		}
		test(value)
	}

}
