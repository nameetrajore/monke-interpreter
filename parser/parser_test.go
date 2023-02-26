package parser

import (
	"testing"
	"monke/ast"
	"monke/lexer"
)

func TestIdentifierExpression(t *testing.T) {
input := "foobar;"
l := lexer.New(input)
p := New(l)
program := p.ParseProgram()
checkParserErrors(t, p)
if len(program.Statements) != 1 {
t.Fatalf("program has not enough statements. got=%d",
len(program.Statements))
}
stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
if !ok {
t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
program.Statements[0])
}
ident, ok := stmt.Expression.(*ast.Identifier)
if !ok {
t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
}
if ident.Value != "foobar" {
t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
}
if ident.TokenLiteral() != "foobar" {
t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
ident.TokenLiteral())
}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0{
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
