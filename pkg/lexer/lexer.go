/*
* Package lexer is used to transform the source code into tokens.
* it uses the original source code representation from the token package
* to transform user generated source code into tokens.
* The lexer is not concerned about the syntax of the source code, but only
* about the token representation.

* Example:
* user input:
* let x = 1 + 2
* output
* [
*     LET,
*     IDENTIFIER("x"),
*     EQUAL_SIGN,
*     INTEGER(1),
*     PLUS_SIGN,
*     INTEGER(2),
*     SEMICOLON,
*     EOF
 */

package lexer

import "github.com/maxwellgithinji/jaba/pkg/token"

// Lexer defines properties required to turn source code into tokens
type Lexer struct {
	// input represent the source code to be tokenized.
	input string

	// position represents the current position in the source code. it points to the to the index of the current character being read.
	position int

	// readPosition represents the next position in the source code. it points to the index of the next character after the position.
	readPosition int

	// ch represents the current character being examined. (Currently only ASCII characters are supported)
	ch byte // TODO: change to rune to support unicode characters
}

// New returns a new lexer for the input.
// It also reads the first character of the input and advances the read position to the next character.
func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

// readChar reads the next character and advances the read position in the input string (source code).
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 is an Ascii code for null
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken returns the next token in the input.
// it converts the input character to a token
// It then advanced the read position so the next call to NextToken will return the next token in the input.
// finally, it returns the token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = newToken(token.PLUS, l.ch)

	case '-':
		tok = newToken(token.MINUS, l.ch)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NEQ,
				Literal: literal,
			}

		} else {
			tok = newToken(token.NOPE, l.ch)
		}

	case '*':
		tok = newToken(token.ASTERISK, l.ch)

	case '/':
		tok = newToken(token.SLASH, l.ch)

	case '<':
		tok = newToken(token.LT, l.ch)

	case '>':
		tok = newToken(token.GT, l.ch)

	case ',':
		tok = newToken(token.COMMA, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)

	case ')':
		tok = newToken(token.RPAREN, l.ch)

	case '{':
		tok = newToken(token.LBRACE, l.ch)

	case '}':
		tok = newToken(token.RBRACE, l.ch)

	case 0:
		tok.Literal = "" // EOF literal is an empty string
		tok = newToken(token.EOF, l.ch)

	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INTEGER
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

// newToken returns a new token with the given type and literal.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// readIdentifier reads an identifier and advances the read position until it encounters a non-letter character.
func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// isLetter returns true if the given character is a letter.
// we also include the underscore character as a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// skipWhitespace skips over all the whitespace characters in the input.
// jaba does not care about the whitespace characters like ruby or python.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads an integer and advances the read position until it encounters a non-digit character.
func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// isDigit returns true if the given character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// peekChar returns the next character in the input without advancing the read position.
// it has the same behavior as the readChar function except that it does not advance the read position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // 0 is an Ascii code for null
	} else {
		return l.input[l.readPosition]
	}
}

// readString loops until it encounters a closing quote or the end of the input and returns the string enclosed by the quotes
func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}
