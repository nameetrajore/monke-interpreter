package parser

import (
	"fmt"
	"monke/ast"
	"monke/lexer"
	"monke/token"
	"strconv"
)

// here we use iota to give the folowing constants incrementing numbers as values.
// The blank identifier _ takes the zero value and the following constants values from 1 to 7.
// The values dont matter but the increasing order of the values do
// This sets the precedence of operators
const (
	_ int = iota
	LOWEST
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT // *
	PREFIX // -X or !X
	CALL //myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser struct contains a pointer to the lexer
// errors to contain the errors we encounter during the parsing of the program
// currToken contains the currentToken in the program
// peekToken contains the nextToken in the program
type Parser struct {
	l *lexer.Lexer
	errors []string
	currToken token.Token
	peekToken token.Token
	// prefixParseFns and infixParseFns maps the type of token to the type of function
	// that should be called when a particular token is encountered.
	// Each token type can have up to two parsing functions associated with it
	// depending on whether the token is found in prefix or infix position.
	prefixParseFns map[token.TokenType]prefixParseFn // maps tokens to functions
	infixParseFns map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression // the argument passed is the 'left side' of the infix operator that's being parsed
)

//TODO: Describe
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

//TODO: Describe
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

//creates a new Parser
// initializes currToken and peekToken
func New(l *lexer.Lexer) *Parser{
	// []string{} is an empty slice
	p := &Parser{ l: l, errors: []string{} }
	//Initializing currToken and peekToken
	p.nextToken() // advances peekToken to the first token
	p.nextToken() // advances currToken to the first token and peekToken to the second token

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// mapping tokens to parsing functions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression{
	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if(err!=nil){
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append (p.errors, msg)
		return nil
	}

	lit.Value = value;

	return lit;
}

func (p *Parser) Errors() []string{
	return p.errors
}

// advances peekToken and currToken by 1
func (p *Parser) nextToken(){
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// We create an empty ast.Program
// We loop through the tokens in the given program
// The meaningful statements are added to ~program~
// To skip the semicolon and to jump to the next statement, we call p.nextToken()
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// On the basis of the starting token it is chosen what kind of a statement we have,
// and the necessary parsing for it
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// XXX: why does it return a pointer and not ast.Statement directly like parseStatement?
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	// if peekToken does not have an identifier, the statement is an invalid let statement, therfore return nil
	if !p.expectPeek(token.IDENT){
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	// if peekToken does not have '=', the statement is an invalid let statement, therefore return nil
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Skipping expressions for now
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.currToken}

	// TODO: Skipping expressions for now
	p.nextToken()
	for p.currTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p. nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
		Left: left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// to check if p.currToken is same as the token we expect
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// helped function for expectPeek()
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek() checks if the next token is the one we expect
// If yes, it returns true and the program continues
// If no, we return false and an Error Message is sent
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p. peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Error function
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Error handling for prefix expressions
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t);
	p.errors = append(p.errors, msg);
}
