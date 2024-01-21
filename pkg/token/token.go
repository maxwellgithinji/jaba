/*
* Package token is used to represent original source code.
* it includes the token type and the token value
 */
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
	// ILLEGAL represents a token that we don't recognize.
	ILLEGAL TokenType = "ILLEGAL"

	// EOF represents the end of the file. It helps the parser to know when to stop parsing.
	EOF TokenType = "EOF"

	// Identifier + Literals

	// IDENTIFIER represents names given by the user to variables and functions. e.g. foo, bar x, y, z
	IDENTIFIER TokenType = "IDENTIFIER"

	// INTEGER represents the number values e.g 1, 2, 3
	INTEGER TokenType = "INTEGER"

	// Operations

	// ASSIGN represents the assignment operation. eg. x = 1
	ASSIGN TokenType = "="

	// PLUS represents the addition operation.
	PLUS TokenType = "+"

	// Delimiters (Special Characters)

	// COMMA represents the comma operator.
	COMMA TokenType = ","

	// SEMICOLON represents the semicolon operator.
	SEMICOLON TokenType = ";"

	// LPAREN represents the left parenthesis operator.
	LPAREN TokenType = "("

	// RPAREN represents the right parenthesis operator.
	RPAREN TokenType = ")"

	// LBRACE represents the left brace operator.
	LBRACE TokenType = "{"

	// RBRACE represents the right brace operator.
	RBRACE TokenType = "}"

	// 	Keywords (Are reserved for the language and cannot be used as identifiers)

	// FUNCTION represents the keyword function.
	FUNCTION TokenType = "FUNCTION"

	// LET represents the keyword let. it is used to declare variables.
	LET TokenType = "LET"
)
