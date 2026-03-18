package Ast

// ==========================================
// STATEMENTS
// ==========================================

// Statement — маркерный интерфейс для всех инструкций.
// (аналог abstract class Statement в C#)
type Statement interface {
	stmtNode()
}

// ExpressionStatement: инструкция-обёртка над выражением
type ExpressionStatement struct {
	Expr Expression
}

func (*ExpressionStatement) stmtNode() {}

// PrintStatement: print x + 2;
type PrintStatement struct {
	Expr Expression
}

func (*PrintStatement) stmtNode() {}

// VarStatement: var x = 5;
type VarStatement struct {
	Name        string
	Initializer Expression // nil допустим, как и null в C#
}

func (*VarStatement) stmtNode() {}

// BlockStatement: блок { ... }
type BlockStatement struct {
	Statements []Statement
}

func (*BlockStatement) stmtNode() {}

// IfStatement: if (...) ... else ...
type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement // nil если else нет
}

func (*IfStatement) stmtNode() {}

// WhileStatement: while (...) ...
type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (*WhileStatement) stmtNode() {}
