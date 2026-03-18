package Ast

import "CompilatorOnGo/Core/Lexer"

type Expression interface {
	exprNode()
}

type NumberExpression struct {
	Value float64
}

func (*NumberExpression) exprNode() {}

type StringExpression struct {
	Value string
}

func (*StringExpression) exprNode() {}

type VariableExpression struct {
	Name string
}

func (*VariableExpression) exprNode() {}

type BinaryExpression struct {
	Left     Expression
	Operator Lexer.TokenType
	Right    Expression
}

func (*BinaryExpression) exprNode() {}

type UnaryExpression struct {
	Operator Lexer.TokenType
	Right    Expression
}

func (*UnaryExpression) exprNode() {}

type AssignExpression struct {
	Name  string
	Value Expression
}

func (*AssignExpression) exprNode() {}
