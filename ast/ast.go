package ast

import (
	"bytes"

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

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64       // the value of the integer
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
