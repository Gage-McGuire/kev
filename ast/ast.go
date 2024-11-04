package ast

import (
	"bytes"
	"strings"

	"github.com/kev/token"
)

// Node is the interface that all nodes in the AST implement
// holds the literal value of the token
// and a string representation of the node
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is the interface that all statement nodes in the AST implement
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface that all expression nodes in the AST implement
type Expression interface {
	Node
	expressionNode()
}

// Holds the root node of every AST that the parser produces
type Program struct {
	Statements []Statement
}

// Represents a var statement
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier // variable name
	Value Expression  // variable value being set
}

// Represents an identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string      // the name/value of the identifier
}

// Represents an return statement
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression  // the value being returned
}

// Represents an expression statement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

// Represents an integer literal
type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64       // the value of the integer
}

// Represents a prefix expression with a prefix operator
type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string      // the operator, e.g. !
	Right    Expression  // the right expression
}

// Represents an infix expression with an infix operator
type InfixExpression struct {
	Token    token.Token // the infix token, e.g. +
	Left     Expression  // the left expression
	Operator string      // the operator, e.g. +
	Right    Expression  // the right expression
}

// Represents a boolean expression
type Boolean struct {
	Token token.Token
	Value bool
}

// Represents a if expression
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// Represents a block statement
type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

// Represents a function literal
type FunctionLiteral struct {
	Token      token.Token // the 'func' token
	Parameters []*Identifier
	Body       *BlockStatement
}

// Represents a call expression
type CallExpression struct {
	Token     token.Token  // the '(' token
	Function  Expression   // the function being called
	Arguments []Expression // the arguments being passed to the function
}

// variable
func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

// identifier
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// return
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// expression
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// integer
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// prefix
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// infix
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// boolean
func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// if
func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// block
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// function
func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// call
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// gets the root node of the AST
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// converts the AST to a string
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// converts the var statement to a string
func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// converts the return statement to a string
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// converts the expression statement to a string
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// converts the identifier to a string
func (i *Identifier) String() string {
	return i.Value
}

// converts the integer literal to a string
func (il *IntegerLiteral) String() string {
	return il.TokenLiteral()
}

// converts the boolean to a string
func (b *Boolean) String() string {
	return b.TokenLiteral()
}

// converts the prefix expression to a string
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// converts the infix expression to a string
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// converts the if expression to a string
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// converts the block statement to a string
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// converts the function literal to a string
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

// converts the call expression to a string
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
