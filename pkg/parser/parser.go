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

// Parser defines properties requires for parsing and turning tokens to AST nodes
type Parser struct {
	// l is a pointer to an instance of lexer. when used, it calls the next token with its New() method
	l *lexer.Lexer

	//currentToken holds the value of the current token under examination
	currentToken token.Token

	// peekToken holds the value of the next token
	peekToken token.Token

	// errors holds a list of errors that occur when parsing
	errors []string

	// prefixParseFns holds a map of prefix functions
	prefixParseFns map[token.TokenType]prefixParseFn

	// infixParseFns holds a map of infix functions
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
	p.registerPrefix(token.NOPE, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

// precedences is a hashmap containing infix operator tokens mapped to respective precedence values
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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
	// Visualizes parseExpressionStatement
	defer untrace(trace("parseExpressionStatement"))

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
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Visualizes parseExpression
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseError(p.currentToken.Type)
		return nil
	}

	leftExpression := prefix()

	// the loop helps the parser find the whole expression
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		// return the left expression if no infix is found
		if infix == nil {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

// noPrefixParseError returns a formatted error when parser encounters no prefix
func (p *Parser) noPrefixParseError(tokenType token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	p.errors = append(p.errors, message)
}

// parseIdentifier returns a representation of an identifier  which contains the token as sIDENTIFIER and the value
// Note: we can return ast.Identifier struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

// parseIntegerLiteral returns a representation of an integer literal which contains the token and value in int64 format
// Note: we can return ast.IntegerLiteral struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// Visualizes parseIntegerLiteral
	defer untrace(trace("parseIntegerLiteral"))

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

// parsePrefixExpression returns a representation of a prefix expression which contains an expression on the left and right side
// Note: we can return ast.IntegerLiteral struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parsePrefixExpression() ast.Expression {
	// Visualizes parsePrefixExpression
	defer untrace(trace("parsePrefixExpression"))

	// left side
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	// skip to the right side
	p.nextToken()

	// parse the expression on the right side
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// peekPrecedence returns the precedence associated with the peek token
// If the peek token has no precedence, it defaults to LOWEST.
func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

// currentPrecedence returns the precedence associated with the current token
// If the current token has no precedence, it defaults to LOWEST.
func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

// parseInfixExpression returns a representation of an infix operator that contains the left expression, operator and right expression
// Note: we can return ast.InfixExpression struct since it fulfills ast.Expression interface by implementing its methods
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Visualizes parseInfixExpression
	defer untrace(trace("parseInfixExpression"))

	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedences := p.currentPrecedence()

	// skip to the right side
	p.nextToken()

	// parse the expression on the right side
	expression.Right = p.parseExpression(precedences)

	return expression
}

// parseBoolean uses go boolean syntax to parse the value of the expression
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIS(token.TRUE)}
}

// parseGroupedExpression uses the left parenthesis to parse set the precedence
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

// parseIfExpression returns a block statement node with the parsed expression
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parseBlockStatement returns a node representing a block statement.
// it parses the block until it encounters } which signifies end of block
// or if it encounters an EOF
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}

	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIS(token.RBRACE) && !p.currentTokenIS(token.EOF) {
		statement := p.parseStatement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

// parseFunctionLiteral returns a node representing a function literal
func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

// parseFunctionParameters returns a list of identifiers that represent function parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// allow empty parameters
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, identifier)

	// parse function parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression returns a node that represents the function call expression
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Function: function}

	expression.Arguments = p.parseCallArguments()

	return expression
}

// parseCallArguments is a helper function that parses the arguments of a function call
func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	// allow empty arguments
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return arguments
	}

	p.nextToken()

	arguments = append(arguments, p.parseExpression(LOWEST))

	// parse function parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return arguments
}
