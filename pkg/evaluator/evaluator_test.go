package evaluator

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/lexer"
	"github.com/maxwellgithinji/jaba/pkg/object"
	"github.com/maxwellgithinji/jaba/pkg/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
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

func testNullObject(t *testing.T, object object.Object) bool {
	if object != NULL {
		t.Errorf("object is not NULL, got %T", object)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 -10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
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
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10};", 10},
		{"if(false) {10};", nil},
		{"if (1) {10};", 10},
		{"if (1 < 2) {10};", 10},
		{"if (1 > 2) {10};", nil},
		{"if (1 > 2) {10} else {20};", 20},
		{"if (1 < 2) {10} else {20};", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1){
					return 10;
				}
				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"9; return 10;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; 9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"unknown operation: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) {true + false}",
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1){
                    return true + false;
                }
                return 1;
			}
			`,
			"unknown operation: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
		{
			`"hello" - "world"`,
			"unknown operation: STRING - STRING",
		},
	}

	for _, tt := range test {
		evaluated := testEval(tt.input)

		errorObject, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("evaluated is not *object.Error, got: %T(%+v)", evaluated, evaluated)
			continue
		}
		if errorObject.Message != tt.expected {
			t.Errorf("errorObject.Message is not %s, got %s", tt.expected, errorObject.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x + 2};"

	evaluated := testEval(input)

	function, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("evaluated is not *object.Function, got: %T(%+v)", evaluated, evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("len(function.Parameters) is not 1, got: %d", len(function.Parameters))
	}

	if function.Parameters[0].String() != "x" {
		t.Fatalf("function.Parameters[0] is not 'x', got: %q", function.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if function.Body.String() != expectedBody {
		t.Fatalf("function.Body is not %q, got: %q", expectedBody, function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) {x}; identity(5);", 5},
		{"let identity = fn(x) {return x}; identity(5);", 5},
		{"let double = fn(x) {return x * 2}; double(5);", 10},
		{"let add = fn(x, y) {return x + y}; add(5, 5);", 10},
		{"let add = fn(x, y) {return x + y}; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) {x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`

	evaluated := testEval(input)

	stringObject, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("evaluated is not *object.String, got: %T(%+v)", evaluated, evaluated)
	}

	if stringObject.Value != "hello world" {
		t.Fatalf("stringObject.Value is not %q, got %q", input, stringObject.Value)
	}
}

func TestStringConcatenation(t *testing.T) {

	input := `"hello" + " " + "world";`

	evaluated := testEval(input)

	stringObject, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("evaluated is not *object.String, got: %T(%+v)", evaluated, evaluated)
	}

	if stringObject.Value != "hello world" {
		t.Fatalf("stringObject.Value is not %q, got %q", "hello world", stringObject.Value)
	}

}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("");`, 0},
		{`len("four");`, 4},
		{`len("hello world")`, 11},
		{`len(1);`, "argument to len not supported, got: INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got: 2 want: 1"},
		{`len([1, 2, 3]);`, 3},
		{`len([]);`, 0},
		{`first([1, 2, 3])`, 1},
		{`first(1)`, "argument to first must be an array, got: INTEGER"},
		{`first([])`, nil},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument to last must be an array, got: INTEGER"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "argument to push must be an array, got: INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))

		case string:
			errorObject, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("evaluated is not *object.Error, got: %T(%+v)", evaluated, evaluated)
				continue
			}
			if errorObject.Message != expected {
				t.Errorf("errorObject.Message is not %s, got %s", expected, errorObject.Message)
			}

		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("evaluated is not *object.Array, got: %T(%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Fatalf("len(array.Elements) is not %d, got: %d", len(expected), len(array.Elements))
				continue
			}

			for i, element := range array.Elements {
				testIntegerObject(t, element, int64(expected[i]))
			}

		default:

		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	evaluated := testEval(input)

	arrayObject, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("evaluated is not *object.Array, got: %T(%+v)", evaluated, evaluated)
	}

	if len(arrayObject.Elements) != 3 {
		t.Fatalf("len(arrayObject.Elements) is not 3, got: %d", len(arrayObject.Elements))
	}

	testIntegerObject(t, arrayObject.Elements[0], 1)
	testIntegerObject(t, arrayObject.Elements[1], 4)
	testIntegerObject(t, arrayObject.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i]", 1},
		{"[1, 2, 3][1 + 1]", 3},
		{"let myArray = [1, 2, 3]; myArray[2]", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}
}
