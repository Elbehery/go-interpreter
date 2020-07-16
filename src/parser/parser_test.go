package parser

import (
	"testing"

	"../ast"
	"../lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
						let x = 5;
						let y = 10;
						let foobar = 838383; `
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("Parser had %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("Parser Error : %q ", msg)
	}

	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
						return 5;
						return 2;
						return -19;
					`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	exp := 3
	if len(program.Statements) != exp {
		t.Fatalf("program statement contains %d , but expected %d", len(program.Statements), exp)
	}

	for _, stmt := range program.Statements {
		rtrn, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt Type is not an *ast.ReturnStatement, but %T", stmt)
			continue
		}

		if rtrn.TokenLiteral() != "return" {
			t.Errorf("stmt Value is not 'return', but %v", rtrn.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements size is wrong, got  %v", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statment type is not ExpressionStatement, got %T", program.Statements[0])
	}
	id, ok := stmt.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("ExpressionStatement value is not Identifier, got %T", stmt.Value)
	}
	if id.Value != "foobar" {
		t.Fatalf("identifier's Value is wrong, expected: %v, but got: %v", "foobar", id.Value)
	}
	if id.TokenLiteral() != "foobar" {
		t.Fatalf("identifier's TokenLiteral() is wrong, expected: %v, but got: %v", "foobar", id.Value)
	}
}

func TestIntegralExpression(t *testing.T) {
	input := "1;"
	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program length is wrong, expected %v, got %v", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statment[0] type is wrong. expected: "+
			"*ast.ExpressionStatement, got: %T", program.Statements[0])
	}
	val, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression type is wrong. expected IntegerLiteral, got %T", stmt.Value)
	}
	if val.Value != 1 {
		t.Fatalf("expression's value is wrong, expected 1, got %v", val.Value)
	}
	if val.TokenLiteral() != "1" {
		t.Fatalf("expression's TokenLiteral() is wrong, expected 1, got %v", val.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program's number of statements is wrong! expected: 1, got %v", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.statement[0] type is wrong. expected *ast.ExpressionStatement, got %T", program.Statements[0])
		}
		prefix, ok := stmt.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("program's expression type is wrong. expected PrefixExpression, got %T", stmt.Value)
		}
		if prefix.Operator != tt.operator {
			t.Fatalf("prefix expression's operator parsing is wrong."+
				"expected %v, got %v", tt.operator, prefix.Operator)
		}

		// test right expression
		iVal, ok := prefix.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("prefix expression's operand is wrong."+
				"expected *ast.IntegerLiteral, got %T", prefix.Right)
		}
		if iVal.Value != tt.value {
			t.Fatalf("prefix expression's operand is wrong."+
				"expected %v, got %v", tt.value, iVal.Value)
		}
	}
}
