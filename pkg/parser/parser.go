/*
* Package parser is responsible for transforming tokens into an abstract syntax tree (AST)
* The parser will be a recursive descent parser, in particular the top down operator precedence parser (Pratt Parser).
 */
package parser

import (
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
}

// New returns a new Parser. it also reads 2 tokens to initialize the current and peek tokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	default:
		return nil
	}
}

// parseLetStatement parses a "let" statement and returns its AST representation
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
		return false
	}
}

// peekTokenIs returns true if the next token is the given type
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}
