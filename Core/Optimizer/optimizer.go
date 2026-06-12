package Optimizer

import (
	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser/Ast"
)

// Optimizer walks the AST and replaces constant sub-expressions with their
// computed results (constant folding). The semantics of the program are
// guaranteed to be unchanged: division by zero is intentionally left unfolded
// so the interpreter still produces the correct runtime error.
type Optimizer struct{}

func New() *Optimizer {
	return &Optimizer{}
}

// Optimize returns a new, semantically equivalent AST with constant
// sub-expressions replaced by literal nodes.
func (o *Optimizer) Optimize(statements []Ast.Statement) []Ast.Statement {
	result := make([]Ast.Statement, len(statements))
	for i, stmt := range statements {
		result[i] = o.optimizeStatement(stmt)
	}
	return result
}

func (o *Optimizer) optimizeStatement(stmt Ast.Statement) Ast.Statement {
	switch s := stmt.(type) {
	case *Ast.ExpressionStatement:
		return &Ast.ExpressionStatement{Expr: o.optimizeExpr(s.Expr)}

	case *Ast.PrintStatement:
		return &Ast.PrintStatement{Expr: o.optimizeExpr(s.Expr)}

	case *Ast.VarStatement:
		init := s.Initializer
		if init != nil {
			init = o.optimizeExpr(init)
		}
		return &Ast.VarStatement{Name: s.Name, Initializer: init}

	case *Ast.FunctionStatement:
		body := make([]Ast.Statement, len(s.Body))
		for i, bodyStmt := range s.Body {
			body[i] = o.optimizeStatement(bodyStmt)
		}
		return &Ast.FunctionStatement{Name: s.Name, Params: s.Params, Body: body}

	case *Ast.BlockStatement:
		stmts := make([]Ast.Statement, len(s.Statements))
		for i, inner := range s.Statements {
			stmts[i] = o.optimizeStatement(inner)
		}
		return &Ast.BlockStatement{Statements: stmts}

	case *Ast.IfStatement:
		cond := o.optimizeExpr(s.Condition)
		then := o.optimizeStatement(s.ThenBranch)
		var elseBranch Ast.Statement
		if s.ElseBranch != nil {
			elseBranch = o.optimizeStatement(s.ElseBranch)
		}
		return &Ast.IfStatement{Condition: cond, ThenBranch: then, ElseBranch: elseBranch}

	case *Ast.WhileStatement:
		return &Ast.WhileStatement{
			Condition: o.optimizeExpr(s.Condition),
			Body:      o.optimizeStatement(s.Body),
		}

	case *Ast.ReturnStatement:
		val := s.Value
		if val != nil {
			val = o.optimizeExpr(val)
		}
		return &Ast.ReturnStatement{Value: val}

	default:
		return stmt
	}
}

func (o *Optimizer) optimizeExpr(expr Ast.Expression) Ast.Expression {
	switch e := expr.(type) {
	case *Ast.BinaryExpression:
		left := o.optimizeExpr(e.Left)
		right := o.optimizeExpr(e.Right)
		if folded, ok := foldBinary(left, e.Operator, right); ok {
			return folded
		}
		return &Ast.BinaryExpression{Left: left, Operator: e.Operator, Right: right}

	case *Ast.UnaryExpression:
		right := o.optimizeExpr(e.Right)
		if folded, ok := foldUnary(e.Operator, right); ok {
			return folded
		}
		return &Ast.UnaryExpression{Operator: e.Operator, Right: right}

	case *Ast.AssignExpression:
		return &Ast.AssignExpression{Name: e.Name, Value: o.optimizeExpr(e.Value)}

	case *Ast.CallExpression:
		callee := o.optimizeExpr(e.Callee)
		args := make([]Ast.Expression, len(e.Arguments))
		for i, arg := range e.Arguments {
			args[i] = o.optimizeExpr(arg)
		}
		return &Ast.CallExpression{Callee: callee, Arguments: args}

	default:
		// NumberExpression, StringExpression, BooleanExpression, VariableExpression
		return expr
	}
}

// foldBinary attempts to evaluate a binary operation on two literal nodes.
// Returns (result, true) on success, (nil, false) when folding is not possible.
func foldBinary(left Ast.Expression, op Lexer.TokenType, right Ast.Expression) (Ast.Expression, bool) {
	switch op {
	case Lexer.PLUS:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.NumberExpression{Value: l.Value + r.Value}, true
			}
		}
		if l, ok := left.(*Ast.StringExpression); ok {
			if r, ok := right.(*Ast.StringExpression); ok {
				return &Ast.StringExpression{Value: l.Value + r.Value}, true
			}
		}

	case Lexer.MINUS:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.NumberExpression{Value: l.Value - r.Value}, true
			}
		}

	case Lexer.STAR:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.NumberExpression{Value: l.Value * r.Value}, true
			}
		}

	case Lexer.SLASH:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				// Do NOT fold division by zero: preserve the runtime error.
				if r.Value == 0 {
					return nil, false
				}
				return &Ast.NumberExpression{Value: l.Value / r.Value}, true
			}
		}

	case Lexer.GT:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value > r.Value}, true
			}
		}

	case Lexer.GTEQ:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value >= r.Value}, true
			}
		}

	case Lexer.LT:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value < r.Value}, true
			}
		}

	case Lexer.LTEQ:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value <= r.Value}, true
			}
		}

	case Lexer.EQEQ:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value == r.Value}, true
			}
		}
		if l, ok := left.(*Ast.StringExpression); ok {
			if r, ok := right.(*Ast.StringExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value == r.Value}, true
			}
		}
		if l, ok := left.(*Ast.BooleanExpression); ok {
			if r, ok := right.(*Ast.BooleanExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value == r.Value}, true
			}
		}

	case Lexer.NEQ:
		if l, ok := left.(*Ast.NumberExpression); ok {
			if r, ok := right.(*Ast.NumberExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value != r.Value}, true
			}
		}
		if l, ok := left.(*Ast.StringExpression); ok {
			if r, ok := right.(*Ast.StringExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value != r.Value}, true
			}
		}
		if l, ok := left.(*Ast.BooleanExpression); ok {
			if r, ok := right.(*Ast.BooleanExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value != r.Value}, true
			}
		}

	case Lexer.AND:
		if l, ok := left.(*Ast.BooleanExpression); ok {
			if r, ok := right.(*Ast.BooleanExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value && r.Value}, true
			}
		}

	case Lexer.OR:
		if l, ok := left.(*Ast.BooleanExpression); ok {
			if r, ok := right.(*Ast.BooleanExpression); ok {
				return &Ast.BooleanExpression{Value: l.Value || r.Value}, true
			}
		}
	}

	return nil, false
}

// foldUnary attempts to evaluate a unary operation on a literal node.
func foldUnary(op Lexer.TokenType, right Ast.Expression) (Ast.Expression, bool) {
	switch op {
	case Lexer.MINUS:
		if r, ok := right.(*Ast.NumberExpression); ok {
			return &Ast.NumberExpression{Value: -r.Value}, true
		}
	case Lexer.EXCL:
		if r, ok := right.(*Ast.BooleanExpression); ok {
			return &Ast.BooleanExpression{Value: !r.Value}, true
		}
	}
	return nil, false
}
