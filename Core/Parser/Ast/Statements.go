package Ast

type Statement interface {
	stmtNode()
}

type ExpressionStatement struct {
	Expr Expression
}

func (*ExpressionStatement) stmtNode() {}

type PrintStatement struct {
	Expr Expression
}

func (*PrintStatement) stmtNode() {}

type VarStatement struct {
	Name        string
	Initializer Expression
}

func (*VarStatement) stmtNode() {}

type FunctionStatement struct {
	Name   string
	Params []string
	Body   []Statement
}

func (*FunctionStatement) stmtNode() {}

type BlockStatement struct {
	Statements []Statement
}

func (*BlockStatement) stmtNode() {}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (*IfStatement) stmtNode() {}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (*WhileStatement) stmtNode() {}

type ReturnStatement struct {
	Value Expression
}

func (*ReturnStatement) stmtNode() {}
