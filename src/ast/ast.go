package ast

import (
	"bytes"
	"playground/go-interpreter/src/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var buf bytes.Buffer

	for _, stmt := range p.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LetStatement) statementNode() {}
func (l *LetStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(l.Token.Literal + " ")
	buf.WriteString(l.Name.String())
	buf.WriteString(" = ")

	if l.Value != nil {
		buf.WriteString(l.Value.String())
	}

	buf.WriteString(";")
	return buf.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.TokenLiteral()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(rs.Token.Literal + " ")
	if rs.ReturnValue != nil {
		buf.WriteString(rs.ReturnValue.String())
	}

	buf.WriteString(";")
	return buf.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	var buf bytes.Buffer

	buf.WriteString("(")
	buf.WriteString(pe.Operator)
	buf.WriteString(pe.Right.String())
	buf.WriteString(")")

	return buf.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var buf bytes.Buffer

	buf.WriteString("(")
	buf.WriteString(ie.Left.String())
	buf.WriteString(" " + ie.Operator + " ")
	buf.WriteString(ie.Right.String())
	buf.WriteString(")")

	return buf.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.TokenLiteral()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var buf bytes.Buffer

	buf.WriteString("if")
	buf.WriteString(ie.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buf.WriteString("else ")
		buf.WriteString(ie.Alternative.String())
	}
	return buf.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var buf bytes.Buffer
	for _, s := range bs.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var buf bytes.Buffer

	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	buf.WriteString(fl.TokenLiteral())
	buf.WriteString("(")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(") ")
	buf.WriteString(fl.Body.String())

	return buf.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var buf bytes.Buffer

	args := make([]string, 0)
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	buf.WriteString(ce.Function.String())
	buf.WriteString("(")
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteString(")")

	return buf.String()
}
