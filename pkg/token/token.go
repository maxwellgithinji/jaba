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

	// PLUS represents the addition operation. eg. x + 1
	PLUS TokenType = "+"

	// MINUS represents the subtraction operation. eg. x - 1
	MINUS TokenType = "-"

	// NOPE represents the negation operation. eg. !x
	NOPE TokenType = "!"

	// ASTERISK represents the multiplication operation. eg. x * 1
	ASTERISK TokenType = "*"

	// SLASH represents the division operation. eg. x / 1
	SLASH TokenType = "/"

	// LT represents the less than operation. eg. x < 1
	LT TokenType = "<"

	// GT represents the greater than operation. eg. x > 1
	GT TokenType = ">"

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

// keywords defines the language reserves characters that cannot be used as identifiers.
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdentifier returns the token type for the given identifier.
// it also checks if the identifier is a keyword and returns it if so.
func LookupIdentifier(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}

	return IDENTIFIER
}
