/*
* Package ast contains functionality that will help with the construction of
* that data structure that will be used to represent the tokens we pass.
* The data structure in this case is an abstract syntax tree that contains the high level
* details of the token representation.
 */

package ast

import "github.com/maxwellgithinji/jaba/pkg/token"

// Node represents a node in the AST
type Node interface {
	// TokenLiteral returns the actual value of the token
	TokenLiteral() string
}

// Statement represents a structure that does not produce values
type Statement interface {
	// Node ensures each statement implements TokenLiteral
	Node

	// statementNode constructs a statement node
	statementNode()
}

// Expression represents a structure that produces values
type Expression interface {
	// Node ensures each expression implements TokenLiteral
	Node

	// expressionNode constructs ans expression node
	expressionNode()
}

// Program represent the root node of every AST the parser will return
type Program struct {
	// Statements contains a list of statements which are building blocks of a program
	Statements []Statement
}

// TokenLiter returns the root node of the AST
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement should contain the LET token, an identifier and for simplification, an expression.
type LetStatement struct {
	//  Token is token.LET
	Token token.Token

	// The Name is the identifier for binding the expression/statement
	Name *Identifier

	// Value represent both the expression and a statement. e.g. "let x = 5" (statement) and "let y = add(2,2)" (expression) can use the expression as the value
	Value Expression
}

// statementNode represents a node that does not produce values
func (l *LetStatement) statementNode() {}

// TokenLiteral returns the a "let" expression in the AST
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// Identifier contains the structure of an identifier
type Identifier struct {
	// Token is the token.IDENTIFIER
	Token token.Token

	// Value is the actual value the identifier represents
	Value string
}

// expressionNode returns a node that produces a value
func (i *Identifier) expressionNode() {}

// TokenLiteral returns the actual identifier value
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
