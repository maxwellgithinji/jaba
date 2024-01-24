/*
* Package ast contains functionality that will help with the construction of
* that data structure that will be used to represent the tokens we pass.
* The data structure in this case is an abstract syntax tree that contains the high level
* details of the token representation.
 */

package ast

import "github.com/maxwellgithinji/jaba/pkg/token"

// Node represents a single branch in the abstract syntax tree
// It ensures every method that implements it also returns a token literal
type Node interface {
	// TokenLiteral returns the actual value of the token
	TokenLiteral() string
}

// Statement is structure that abstracts a list of tokens that resemble a single statement
// Every method that implements Statement should return a token literal (from Node interface) and construct a
// statement node. Statement does not produce a value. e.g. return 5, return foo
type Statement interface {
	// Node ensures each statement returns a token literal
	Node

	// statementNode constructs a statement node from a combination of tokens that represent a statement
	statementNode()
}

// Expression is a structure that abstracts a list of tokens that represent an expression
// Every method that implements Expression should return a token literal (from Node interface) and construct an
// expression node. Expression produces a value e.g. return add(1, 2), return 1 + 4
type Expression interface {
	// Node ensures each expression returns a token literal
	Node

	// expressionNode constructs an expression node from a combination of tokens that represent an expression
	expressionNode()
}

// Program represents entry point where the root of the AST is initialized and other child nodes are built into the AST
type Program struct {
	// Statements contains a list of nodes (branches) which are building blocks the AST
	Statements []Statement
}

// TokenLiteral method returns the token literal of the first statement in the program.
// (Root node of the AST)
// This method is used by the parser to determine the first token to be executed when a program is run.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement defines the  3 parts of a let let statement: "let", "identifier", "expression"
type LetStatement struct {
	//  Token is token.LET i.e. {Type: LET, literal: "let"}
	Token token.Token

	// The Name is the identifier for binding the expression/statement e.g. {token: IDENTIFIER, value: "foo"}
	Name *Identifier

	// Value represent both the expression and a statement. e.g. "let x = 5" (statement) and "let y = add(2,2)" (expression) can use the expression as the value
	Value Expression
}

// statementNode method is method that is used to indicate that the LetStatement struct
// can be used to construct a statement in the Abstract Syntax Tree (AST).
func (l *LetStatement) statementNode() {}

// TokenLiteral returns the a "let" expression in the AST
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// Identifier represents the 2 parts of an identifier, IDENTIFIER and Value e.g. IDENTIFIER("foo")
type Identifier struct {
	// Token is the token.IDENTIFIER
	Token token.Token

	// Value is the actual value the identifier represents e.g. "foo"
	Value string
}

// expressionNode method is method that is used to indicate that the Identifier struct
// can be used to construct an expression in the Abstract Syntax Tree (AST).
func (i *Identifier) expressionNode() {}

// TokenLiteral returns the actual value of the identifier e.g. returns "foo" in IDENTIFIER("foo")
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// ReturnStatement contains the 2 parts of the return statement, RETURN(expression) e.g. "return add(5,5)"
type ReturnStatement struct {
	// Token is the token.RETURN
	Token token.Token

	// Value is the actual expression being returned e.g. add(5,5), 5, foo, nil. note, we can return both statements and expressions
	Value Expression
}

// statementNode method is method that is used to indicate that the ReturnStatement struct
// can be used to construct a statement in the Abstract Syntax Tree (AST).
func (r *ReturnStatement) statementNode() {}

// TokenLiteral returns an actual value in a return statement e.g. add(5,5), 5, foo, nil
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
