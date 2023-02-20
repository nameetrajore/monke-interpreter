package ast

import "monke/token"


type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface{
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p* Program) TokenLiteral() string {
	if len (p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

// LetStatement is a part of Node and Statement
func (ls *LetStatement) statementNode(){}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

// Identifier is a part of Node and Expression
func (i* Identifier) expressionNode() {}
func (i* Identifier) TokenLiteral() string { return i.Token.Literal }
