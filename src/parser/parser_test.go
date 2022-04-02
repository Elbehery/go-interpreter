package parser

import (
	"playground/go-interpreter/src/ast"
	"playground/go-interpreter/src/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected %v, but got %v statements instead", 3, len(program.Statements))
	}

	testCases := []struct {
		ident string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range testCases {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.ident) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, ident string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	s, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if s.Name.Value != ident {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", ident, s.Name.Value)
		return false
	}

	if s.Name.TokenLiteral() != ident {
		t.Errorf("s.Name not '%s'. got=%s", ident, s.Name)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
   return 5;
   return 10;
   return 993322;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		rs, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if rs.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", rs.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {

	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("expected 1 statement, but got %v instead", len(program.Statements))
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

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %v", msg)
	}
	t.FailNow()
}
