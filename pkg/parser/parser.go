/*
* Package parser is responsible for transforming tokens into an abstract syntax tree (AST)
* The parser will be a recursive descent parser, in particular the top down operator precedence parser (Pratt Parser).
 */
package parser

import (
	"fmt"
	"strconv"

	"github.com/maxwellgithinji/jaba/pkg/ast"
	"github.com/maxwellgithinji/jaba/pkg/lexer"
	"github.com/maxwellgithinji/jaba/pkg/token"
)

// Parser contains a pointer to an instance of lexer, the current and peek token
type Parser struct {
	// l is a pointer to an instance of lexer which repeatedly calls nextToken to read the next token
	l *lexer.Lexer
	//currentToken holds the value of the current token under examination
	currentToken token.Token
	// peekToken holds the value of the next token
	peekToken token.Token

	// errors are returned when the parser encounters a tokens that are not of the expected type
	errors []string

	// prefixParseFns is used to get the correct prefix for the current token
	prefixParseFns map[token.TokenType]prefixParseFn

	// infixParseFns is used to get the correct infix for the current token
	infixParseFns map[token.TokenType]infixParseFn
}

// New returns a new Parser. it also reads 2 tokens to initialize the current and peek tokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)

	p.nextToken()
	p.nextToken()
	return p
}

// nextToken is a helper function to set the current and peek tokens
// The current token is set to the current peek token
// Peek token is set to the next peek token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram returns an AST representing the tokens
// It first constructs the root node of the AST then iterates over the
// tokens while constructing the tree with child nodes and advancing the positions
// until it encounters an EOF therefore returns
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement and returns its AST representation
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement creates an AST representation of a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skip expression until we encounter semicolon

	if !p.currentTokenIS(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// currentTokenIS returns true if the current token is the given type
func (p *Parser) currentTokenIS(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

// expectPeek advances the position to the next token if the given type is true
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	} else {
		p.peekError(tokenType)
		return false
	}
}

// peekTokenIs returns true if the next token is the given type
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

// Errors returns a slice containing all the errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError appends error message to errors when it encounters a peek token that does not match the given type
func (p *Parser) peekError(tokenType token.TokenType) {
	message := fmt.Sprintf("expected next token to be %v, got %v", tokenType, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

// parseReturnStatement creates the AST representation of a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// TODO: We are skipping the expression until we find a semicolon
	if !p.currentTokenIS(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

type (
	// prefixParseFn  parses tokens that are in a prefix position
	prefixParseFn func() ast.Expression

	// infixParseFn parses tokens that are in an infix position
	// The argument passed here is on the left side of the infix operator
	infixParseFn func(ast.Expression) ast.Expression
)

// This iota is used to order the constants based on precedence from the lowest to the highest
const (
	// _ has the value 0
	_ int = iota

	// LOWEST has the value 1
	LOWEST

	// EQUALS has the value 2 (==)
	EQUALS

	// LESSGREATER has the value 3 (< OR >)
	LESSGREATER

	// SUM has the value 4 (+)
	SUM
	// PRODUCT has the value 5 (*)
	PRODUCT

	// PREFIX has the value 6 (-x or !x)
	PREFIX

	// CALL has the value 7. add(x, y)
	CALL
)

// registerPrefix records a prefix token
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix records an infix token
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseExpressionStatement creates the AST representation of an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	// we pass the lowest possible precedence since we are initializing and have nothing to compare against
	statement.Value = p.parseExpression(LOWEST)

	// we wont return an error if the expression in the repl does not end with a semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parseExpression is a helper function to parse supported expressions
// TODO: add more docs to explain supported expressions
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	return leftExpression
}

// parseIdentifier returns the current token and the literal it represents
// Note: we can return ast.Identifier struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

// parseIntegerLiteral returns the current token and the literal it represents
// Note: we can return ast.IntegerLiteral struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	literal.Value = value

	return literal
}
