/*
* Package evaluator uses the object system to evaluate the AST
 */
package evaluator

import (
	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/object"
)

var (
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
