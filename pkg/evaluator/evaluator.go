/*
* Package evaluator uses the object system to evaluate the AST
 */
package evaluator

import (
	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval is a recursive function that that evaluates the AST and returns an object representation as output
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Value)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBooleanToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right) // evaluates expression on the right hand side of the operator
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)   // evaluates expression on the left hand side of the operator
		right := Eval(node.Right) // evaluates expression on the right hand side of the operator
		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

// evalStatements is a helper function that a list of AST statements and returns an object representation as output
func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

// nativeBooleanToBooleanObject is a helper function that converts a native boolean to a boolean object
func nativeBooleanToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

// evalPrefixExpression is a helper function that evaluates a prefix expression, and returns an object representation as output
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNopePrefixOperatorExpression(right)

	case "-":
		return evalMinusPrefixOperatorExpression(right)

	}
	return nil
}

// evalNopeOperatorExpression is a helper function that evaluates a nope operator that appears at the beginning of the expression
// TODO:  should be an error handler instead
func evalNopePrefixOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression is a helper function that evaluates a minus operator that appears at the beginning of the expression
// minus prefix only applies to numbers
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return NULL
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

// evalInfixExpression evaluates an expression that have operands in between themselves
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT: // integer based infix expression
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBooleanToBooleanObject(left == right)

	case operator == "!=":
		return nativeBooleanToBooleanObject(left != right)

	default:
		return NULL
	}
}

// evalIntegerInfixExpression evaluated integer based infix expression
func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}

	case "-":
		return &object.Integer{Value: leftValue - rightValue}

	case "*":
		return &object.Integer{Value: leftValue * rightValue}

	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	case "<":
		return nativeBooleanToBooleanObject(leftValue < rightValue)

	case ">":
		return nativeBooleanToBooleanObject(leftValue > rightValue)

	case "==":
		return nativeBooleanToBooleanObject(leftValue == rightValue)

	case "!=":
		return nativeBooleanToBooleanObject(leftValue != rightValue)

	default:
		return NULL
	}
}
