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
		return evalProgram(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Value)

	case *ast.BlockStatement:
		return evalBlockStatements(node)

	case *ast.ReturnStatement:
		value := Eval(node.Value)
		return &object.ReturnValue{Value: value}

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

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

// evalProgram evaluates the entry point of the program
func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// evalBlockStatements is a helper function that evaluates a list of AST block statements and returns an object representation as output
func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
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

// evalIfExpression an evaluated result of the if expression
func evalIfExpression(i *ast.IfExpression) object.Object {
	condition := Eval(i.Condition)

	if isTruthy(condition) {
		return Eval(i.Consequence)
	} else if i.Alternative != nil {
		return Eval(i.Alternative)
	} else {
		return NULL
	}
}

// isTruthy checks if an expression can be evaluated or skipped
func isTruthy(object object.Object) bool {
	switch object {
	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		return true
	}
}
