package ast

import "monke/token"
import "bytes"
import "strings"

// An interface is a type that defines a set of methods
type Node interface {
	TokenLiteral() string // returns token literal
	String() string // converts the struct to a string format. The way it would appear in the program. For debugging and testing purposes.
 }

// Statements and Expressions are the only 2 kinds of statements we are considering in monke
type Statement interface {
	Node // TokenLiteral()
	statementNode()
}
type Expression interface{
	Node // TokenLiteral()
	expressionNode()
}


type Program struct {
	Statements []Statement
}

// It creates a buffer and writes the return value of each statement's String()
// method to it and returns the buffer as a string
func (p *Program) String() string{
	var out bytes.Buffer //bytes.Buffer creates an empty buffer without any initialization

	for _, s := range p.Statements {
		out.WriteString(s.String()) // writes statements to the buffer
	}

	return out.String() // converts the buffer to a string and returns it
}
func (p *Program) TokenLiteral() string {
	if len (p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement is a part of Node and Statement
type LetStatement struct {
	Token token.Token
	Name *Identifier // XXX: Why is this a pointer?
	Value Expression
}

func (ls *LetStatement) statementNode(){} // for debugging and testing purposes
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer // empty buffer, no initialization needed

	out.WriteString(ls.TokenLiteral() + " ") // appends 'let' BUFFER: 'let'
	out.WriteString(ls.Name.String()) // appends {IDENTIFIER} BUFFER: 'let {IDENTIFIER}'
	out.WriteString(" = ") // appends ' = ' BUFFER: 'let {IDENTIFIER} = '

	// appends the Expression BUFFER: 'let {IDENTIFIER} = {EXPRESSION}'
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";") // appends ';' BUFFER: 'let {IDENTIFIER} = {EXPRESSION};'
	return out.String() // converts buffer to string and returns it
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode(){} // for debugging and testing purposes
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer // empty buffer
	out.WriteString(rs.TokenLiteral() + " ") // appends 'return' BUFFER: 'return'

	// appends {EXPRESSION} BUFFER: 'return {EXPRESSION}'
	if rs.ReturnValue != nil{
		out.WriteString(rs.ReturnValue.String())
	}

	// appends ';' BUFFER: 'return {EXPRESSION};'
	out.WriteString(";")
	return out.String() // converts buffer to string and returns it
}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode(){} // for debugging and testing purposes
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String() // returns the entire expression as it is
	}

	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

// Identifier is an Expression
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value } // returns the identifier as it is


type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {return il.Token.Literal}
func (il *IntegerLiteral) String() string {return il.Token.Literal}

type PrefixExpression struct{
	Token token.Token
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {return pe.Token.Literal};
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer;
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct{
	Token token.Token
	Operator string
	Left Expression
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {return ie.Token.Literal};
func (ie *InfixExpression) String() string {
	var out bytes.Buffer;
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct{
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string{
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token token.Token
	Parameters []*Identifier
	Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
