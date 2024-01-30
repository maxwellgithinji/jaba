package evaluator

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/lexer"
	"github.com/maxwellgithinji/jaba/pkg/object"
	"github.com/maxwellgithinji/jaba/pkg/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNopeOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false}, // Literals are truthy in jaba language
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Fatalf("evaluated is not *object.Integer, got: %T", evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %d, got %d", expected, result.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Fatalf("evaluated is not *object.Boolean, got: %T", evaluated)
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %t, got %t", expected, result.Value)
		return false
	}
	return true
}
