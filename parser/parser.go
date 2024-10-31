package parser

import (
	"strconv"

	"github.com/kev/ast"
	"github.com/kev/lexer"
	token "github.com/kev/token"
)

// Parser struct
type Parser struct {
	l               *lexer.Lexer
	currentToken    token.Token
	peekToken       token.Token
	errors          []string
	prefixParseFunc map[token.TokenType]prefixParseFunc
	infixParseFunc  map[token.TokenType]infixParseFunc
}

type (
	// operators before the expression
	prefixParseFunc func() ast.Expression

	// operators within the expression
	infixParseFunc func(ast.Expression) ast.Expression
)

// precedence levels
const (
	_ int = iota

	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// precedence map
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Initializes a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	// initialize the prefixParseFunc map
	p.prefixParseFunc = make(map[token.TokenType]prefixParseFunc)

	// register the prefix parser functions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	// initialize the infixParseFunc map
	p.infixParseFunc = make(map[token.TokenType]infixParseFunc)

	// register the infix parser functions
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

// Advances the parser by one token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parses all statements in the input
// and adds them to the AST statements array
func (p *Parser) ParseProgram() *ast.Program {
	// construct the root node of the AST
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// iterate through all the tokens in the input
	// and parse them into statements
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// Parses a statement based on the current token
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parses a var statement
func (p *Parser) parseVarStatement() *ast.VarStatement {
	// construct the var statement node
	stmt := &ast.VarStatement{Token: p.currentToken}

	// check if the next token is an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// construct the identifier node
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// check if the next token is an assignment
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: skip the expressions for now

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

// Parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// construct the return statement node
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// move to the next token
	p.nextToken()

	//TODO: skip the expressions for now
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// get the prefix parser function for the current token
	prefix := p.prefixParseFunc[p.currentToken.Type]

	// if the prefix is empty, then skip
	if prefix == nil {
		p.noPrefixParseFuncError(p.currentToken.Type)
		return nil
	}

	// parse the left expression
	leftExp := prefix()

	// loop through the tokens and parse the infix expressions
	// while the next token is not a semicolon and
	// the precedence is less than the next precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFunc[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// Parses a grouped expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	// move to the next token
	p.nextToken()

	// parse the expression within the parentheses
	exp := p.parseExpression(LOWEST)

	// check if the next token is a closing parenthesis
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// Parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// construct the expression statement node
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// parse the expression with the lowest precedence
	stmt.Expression = p.parseExpression(LOWEST)

	// check if the next token is a semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// construct the integer literal node
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	// parse the literal into a int64
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := "could not parse " + p.currentToken.Literal + " as integer"
		p.errors = append(p.errors, msg)
		return nil
	}

	// set the value of the integer literal node
	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// construct the prefix expression node
	// with the current token as the operator
	// this is the prefix operator (e.g. !, -)
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	// move to the next token
	// this will set the right expression
	p.nextToken()

	// parse the right expression
	// with PREFIX precedence
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// construct the infix expression node
	// with the current token as the operator
	// this is the infix operator (e.g. +, -, *, /)
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	// get the precedence of the current token
	precedence := p.currentPrecedence()

	// move to the next token
	// this will set the right expression
	p.nextToken()

	// parse the right expression
	// with the precedence of the current token
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Parses a boolean
func (p *Parser) parseBoolean() ast.Expression {
	// construct the boolean node
	exp := &ast.Boolean{Token: p.currentToken}

	// set the value of the boolean node
	exp.Value = p.currentTokenIs(token.TRUE)

	return exp
}

// Parses an if expression
func (p *Parser) parseIfExpression() ast.Expression {
	// construct the if expression node
	expression := &ast.IfExpression{Token: p.currentToken}

	// check if the next token is a left parenthesis
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// move to the next token
	// this will be the start of the condition
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// check if the next token is a right parenthesis
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// check if the next token is a left brace
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// parse the block statement
	// this will be the consequence
	// parseBlockStatement() will parse until the right brace
	expression.Consequence = p.parseBlockStatement()

	// check if there is an alternative
	if p.peekTokenIs(token.ELSE) {
		// move to the next token
		// this will be the start of the alternative
		p.nextToken()

		// check if the next token is a left brace
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		// parse the block statement
		// this will be the alternative
		// parseBlockStatement() will parse until the right brace
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// Parses a block statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// construct the block statement node
	block := &ast.BlockStatement{Token: p.currentToken}

	// initialize the statements array
	block.Statements = []ast.Statement{}

	// move to the next token
	p.nextToken()

	// parse all the statements within the block
	// until the right brace or end of file
	// parseStatement() will parse each statement and add it to the array
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// Parses an identifier
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

// adds a prefix parser function entry to the prefixParseFunc map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunc) {
	p.prefixParseFunc[tokenType] = fn
}

// adds an infix parser function entry to the infixParseFunc map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunc) {
	p.infixParseFunc[tokenType] = fn
}

// Checks if the current token type is the expected type (t)
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// Checks if the next token type is the expected type (t)
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Uses peekTokenIs() function to check if the next token is of the expected type.
// If it is, it advances the parser to the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Returns the precedence of the next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Returns the errors array
func (p *Parser) Errors() []string {
	return p.errors
}

// Adds an error message to the errors array
// when there is a peek error (the next token is unexpected)
func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be " + string(t) + ", got " + string(p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Adds an error message to the errors array
// when there is no prefix parse function for a token
func (p *Parser) noPrefixParseFuncError(t token.TokenType) {
	msg := "no prefix parse function for " + string(t) + " found"
	p.errors = append(p.errors, msg)
}
