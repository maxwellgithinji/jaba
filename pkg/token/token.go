// Package token is used to represent original source code.
// it includes the token type and the token value
package token

/*
TokenType represents the category of a token.
It is of type string

Pros:
- easy to represent any value type
- easy to debug

Cons:
- not as perfomant as a byte/rune/int
*/
type TokenType string

// Token represents the structure that will hold all the information about a token.
type Token struct {
	// Type defines which category the token belongs to.
	Type TokenType
	// Literal defines the actual value of the token.
	Literal string
}

const (
	// Illegal represents a token that we don't recognize.
	Illegal TokenType = "ILLEGAL"

	// EOF represents the end of the file. It helps the parser to know when to stop parsing.
	EOF TokenType = "EOF"

	// Identifier + Literals

	// Identifier represents names given by the user to variables and functions. e.g. foo, bar x, y, z
	Identifier TokenType = "IDENTIFIER"

	// Integer represents the number values e.g 1, 2, 3
	Integer TokenType = "INTEGER"

	// Operations

	// Assign represents the assignment operation. eg. x = 1
	Assign TokenType = "="

	// Plus represents the addition operation.
	Plus TokenType = "+"

	// Delimiters (Special Characters)

	// Comma represents the comma operator.
	Comma TokenType = ","

	// Semicolon represents the semicolon operator.
	Semicolon TokenType = ";"

	// LParen represents the left parenthesis operator.
	LParen TokenType = "("

	// RParen represents the right parenthesis operator.
	RParen TokenType = ")"

	// LBrace represents the left brace operator.
	LBrace TokenType = "{"

	// RBrace represents the right brace operator.
	RBrace TokenType = "}"

	// LBracket represents the left bracket operator.

	// 	Keywords (Are reserved for the language and cannot be used as identifiers)

	// Function represents the keyword function.
	Function TokenType = "FUNCTION"

	// Let represents the keyword let. it is used to declare variables.
	Let TokenType = "LET"
)
