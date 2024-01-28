/*
* Package ast contains functionality that will help with the construction of
* that data structure that will be used to represent the tokens we pass.
* The data structure in this case is an abstract syntax tree that contains the high level
* details of the token representation.
 */

package ast

import (
	"bytes"
	"strings"

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
	//  Token represents the "let" token
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
	// Token represents the "identifier" token
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
	// Token represent the return token
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

// ExpressionStatement is an expression wrapper that groups expressions
// It fulfils the Statement interface by implementing statementNode() method
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type ExpressionStatement struct {
	// Token represents any token representation being parsed as an expression
	Token token.Token

	// Value is any value representation being parsed as an expression
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
	// Token represent the integer token e.g. "5"
	Token token.Token

	// Value asserts the integer value. e.g. "5" will be returned as 5 of type int64
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

// PrefixExpression represents an expression that is placed on the left side of other expressions e.g ! in !5
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type PrefixExpression struct {
	// Token represent the prefix operator token e.g. !
	Token token.Token

	// Operator is a type of prefix that appears on the left side of the expression e.g. ! in !5
	Operator string

	// Right represents the expression on the right hand side of the operator e.g. 5 in !5
	Right Expression
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the prefix expression
func (p *PrefixExpression) expressionNode() {}

// TokenLiteral returns the actual value of the prefix expression e.g.!5
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String returns a string representation of a PrefixExpression node
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents the 3 parts of an infix expression, left operand, operator, right operand. e.g. 5 + 7
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type InfixExpression struct {
	// Token represent the infix operator token e.g. +
	Token token.Token

	// Left represents the expression on the left hand side of the operator e.g. 5 in 5 + 7
	Left Expression

	// Operator is a type of operator that appears on the left side of the expression e.g. + in 5 + 7
	Operator string

	// Right represents the expression on the right hand side of the operator e.g. 7 in 5 + 7
	Right Expression
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the infix expression
func (i *InfixExpression) expressionNode() {}

// TokenLiteral returns the actual value of the infix expression e.g. 5 + 7
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of an InfixExpression node
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean represents whose value is true or false
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type Boolean struct {
	// Token represents either token.True or token.False
	Token token.Token

	// Value is true or false
	Value bool
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the boolean
func (b Boolean) expressionNode() {}

// TokenLiteral returns the string value of the boolean e.g. "true" or "false"
func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String returns a string representation of a Boolean node
func (b Boolean) String() string {
	return b.Token.Literal
}

// IfExpression represents the composition of an if expression that represents an if-else-condition
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type IfExpression struct {
	// Token represents the if token
	Token token.Token

	// Condition represents the expression the if expression is checking
	Condition Expression

	// Consequence represents the block statement to be executed when the condition is met
	Consequence *BlockStatement

	// Alternative represents the block statement to be executed when the condition is not met (ELSE)
	Alternative *BlockStatement
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the if expression
func (i *IfExpression) expressionNode() {}

// TokenLiteral returns the actual value of the if expression e.g. 5 < 7
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of an IfExpression node
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

// BlockStatement represents a list of statements that can be structured in a block like manner
// It fulfils the Statement interface by implementing statementNode() method
// It by extension fulfills the Node interface which is part of the Statement interface
// by implementing TokenLiteral() and String() methods from the Node interface
type BlockStatement struct {
	// Token represents { which indicated the start of a block statement i.e token.RBRACE {
	Token token.Token

	// Statements represents the list of statements in the block
	Statements []Statement
}

// statementNode method constructs a statement node in the Abstract Syntax Tree (AST) from the block statement
func (b *BlockStatement) statementNode() {}

// TokenLiteral returns the actual value of the block statement
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

// String returns a string representation of a BlockStatement node
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// FunctionLiteral defines the structure of a function which includes the fn token, parameters and the body
// It fulfils the Expression interface by implementing expressionNode() method
// It by extension fulfills the Node interface which is part of the Expression interface
// by implementing TokenLiteral() and String() methods from the Node interface
type FunctionLiteral struct {
	// Token represents the fn token
	Token token.Token

	// Parameters represents the parameters of the function
	Parameters []*Identifier

	// Body represents the body of the function
	Body *BlockStatement
}

// expressionNode method constructs an expression node in the Abstract Syntax Tree (AST) from the function literal
func (f *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the actual value of the function literal
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// String returns a string representation of a FunctionLiteral node
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(f.TokenLiteral())

	out.WriteString("(")

	out.WriteString(strings.Join(params, ", "))

	out.WriteString(") ")

	out.WriteString(f.Body.String())

	return out.String()
}
