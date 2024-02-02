/*
* Package evaluator uses the object system to evaluate the AST
 */
package evaluator

import (
	"fmt"

	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval is a recursive function that that evaluates the AST and returns an object representation as output
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Value, env)

	case *ast.BlockStatement:
		return evalBlockStatements(node, env)

	case *ast.ReturnStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}

	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBooleanToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env) // evaluates expression on the right hand side of the operator
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env) // evaluates expression on the left hand side of the operator
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env) // evaluates expression on the right hand side of the operator
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunctions(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// Identifier
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}

// evalProgram evaluates the entry point of the program
func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch r := result.(type) {

		case *object.ReturnValue:
			return r.Value

		case *object.Error:
			return r
		}
	}

	return result
}

// evalBlockStatements is a helper function that evaluates a list of AST block statements and returns an object representation as output
func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJECT || resultType == object.ERROR_OBJECT {
				return result
			}
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
	return newError("unknown operation: %s %s", operator, right.Type())
}

// evalNopeOperatorExpression is a helper function that evaluates a nope operator that appears at the beginning of the expression
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
		return newError("unknown operation: -%s", right.Type())
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

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operation: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression returns evaluated integer based infix expression
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
		return newError("unknown operation %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIfExpression returns an evaluated result of the if expression
func evalIfExpression(i *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(i.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(i.Consequence, env)
	} else if i.Alternative != nil {
		return Eval(i.Alternative, env)
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

// newError returns a meaningful error message to the user of the jaba program when they write unexpected jaba code
// it uses the standard golang Sprintf to format the error message
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError is a helper function that helps check error early and allows them not to stray far away from their origin
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}
	return false
}

// evalIdentifier uses the environment to get the identifier object otherwise returns an error
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	key, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return key
}

// evalExpressions is a helper function that helps evaluate a list of expressions
// the expressions are evaluated from left to right
func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var evaluated []object.Object

	for _, expression := range expressions {
		result := Eval(expression, env)
		if isError(result) {
			return []object.Object{result}
		}
		evaluated = append(evaluated, result)
	}

	return evaluated
}

// applyFunctions is a helper function that helps evaluate a function considering its scope
// it supports higher order functions (functions that return other functions or pass them as arguments)
// and closures (function that close over the environment they were defined in).
func applyFunctions(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)

	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

// extendFunctionEnv is a helper function that helps extend the environment of a function
// by scoping the function environment in an enclosed hash
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

// unwrapReturnValue is a helper function that helps give the value the function returns after executing
func unwrapReturnValue(result object.Object) object.Object {
	if returnValue, ok := result.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return result
}
