package ast

import "monke/token"
import "bytes"

// any type that also wants to be a "Node" should implement all the functions inside the interface
type Node interface {
	TokenLiteral() string
	String() string // for debugging and testing purposes
}
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
		out.WriteString(s.String())
	}

	return out.String()
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
	Name *Identifier //QUESTION: Why is this a pointer?
	Value Expression
}

func (ls *LetStatement) statementNode(){} // for debugging and testing purposes
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode(){} // for debugging and testing purposes
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil{
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode(){} // for debugging and testing purposes
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

// Identifier is a part of Node and Expression
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }
