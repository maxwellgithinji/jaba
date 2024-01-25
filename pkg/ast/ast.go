/*
* Package ast contains functionality that will help with the construction of
* that data structure that will be used to represent the tokens we pass.
* The data structure in this case is an abstract syntax tree that contains the high level
* details of the token representation.
 */

package ast

import (
	"bytes"

	"github.com/maxwellgithinji/jaba/pkg/token"
)

// Node represents a single branch in the abstract syntax tree
// The implementor fulfills the Node interface by implementing
// TokenLiteral() and String() methods
type Node interface {
	// TokenLiteral returns the actual value of the token
	TokenLiteral() string

	// String returns a string representation of an AST node
	String() string
}

// Statement is structure that abstracts a list of tokens that resemble a single statement
// The implementor fulfills the Statement Interface by implementing the statementNode() method
// and by extension, the Node interface by implementing TokenLiteral() and String() methods
type Statement interface {
	// Node ensures each statement returns a token literal and a debug string
	Node

	// statementNode constructs a statement node from a combination of tokens that represent a statement
	statementNode()
}

// Expression is a structure that abstracts a list of tokens that represent an expression
// The implementor fulfills the Expression Interface by implementing the expressionNode() method
// and by extension, the Node interface by implementing TokenLiteral() and String() methods
type Expression interface {
	// Node ensures each statement returns a token literal and a debug string
	Node

	// expressionNode constructs an expression node from a combination of tokens that represent an expression
	expressionNode()
}

// Program represents entry point where the root of the AST is initialized and other child nodes are built into the AST
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type Program struct {
	// Statements contains a list of nodes (branches) which are building blocks of the AST
	Statements []Statement
}

// TokenLiteral method returns the token literal of the first statement in the program.(Root node of the AST)
// This method is used by the parser to determine the first token to be executed when a program is run.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String returns a string representation of a Program node
func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

// LetStatement defines the  3 parts of a let let statement: "let", "identifier", "expression"
// It fulfils the Statement interface by implementing statementNode() method
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type LetStatement struct {
	//  Token is token.LET i.e. {Type: LET, literal: "let"}
	Token token.Token

	// The Name is the identifier for binding the expression/statement e.g. {token: IDENTIFIER, value: "foo"}
	Name *Identifier

	// Value represent both the expression ("add(2,2)") and a statement ("let x = 5"). statement is already represented by the expression
	Value Expression
}

// statementNode method constructs a statement node in the Abstract Syntax Tree (AST) from the let statement
func (l *LetStatement) statementNode() {}

// TokenLiteral returns the a "let" expression in the AST
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// String returns a string representation of a LetStatement node
func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Identifier represents the 2 parts of an identifier, IDENTIFIER and Value e.g. IDENTIFIER("foo")
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type Identifier struct {
	// Token is the token.IDENTIFIER
	Token token.Token

	// Value is the actual value the identifier represents e.g. "foo"
	Value string
}

// expressionNode method constructs a statement node in the Abstract Syntax Tree (AST) from the identifier
func (i *Identifier) expressionNode() {}

// TokenLiteral returns the actual value of the identifier e.g. returns "foo" in IDENTIFIER("foo")
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of an Identifier value
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement contains the 2 parts of the return statement, RETURN(expression) e.g. "return add(5,5)"
// It fulfils the Statement interface by implementing statementNode() method
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type ReturnStatement struct {
	// Token is the token.RETURN
	Token token.Token

	// Value is the actual expression being returned e.g. add(5,5), 5, foo, nil. note, we can return both statements and expressions
	Value Expression
}

// statementNode method constructs a statement node in the Abstract Syntax Tree (AST) from the return statement
func (r *ReturnStatement) statementNode() {}

// TokenLiteral returns an actual value in a return statement e.g. add(5,5), 5, foo, nil
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// String returns a string representation of a ReturnStatement node
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is an expression wrapper that contains the initial token of the expression and the rest of the expression
// It fulfils the Statement interface by implementing statementNode() method
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type ExpressionStatement struct {
	// Token is the first token of the expression
	Token token.Token

	// Value is the rest of the expression
	Value Expression
}

// statementNode method constructs a statement node in the Abstract Syntax Tree (AST) from the expression statement
func (e *ExpressionStatement) statementNode() {}

// TokenLiteral returns the actual value of the expression e.g. add(5,5), --4, z == a
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// String returns a string representation of an ExpressionStatement node
func (e *ExpressionStatement) String() string {
	if e.Value != nil {
		return e.Value.String()
	}

	return ""
}

// IntegerLiteral represents an integer literal in int64 format
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the integer literal
func (n *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the actual value of the literal in string format e.g. "5"
func (n *IntegerLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// String returns a string representation of an integer literal node
func (n *IntegerLiteral) String() string {
	return n.Token.Literal
}
