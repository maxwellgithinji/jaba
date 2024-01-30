/*
* Package evaluator uses the object system to evaluate the AST
 */
package evaluator

import (
	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/object"
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
