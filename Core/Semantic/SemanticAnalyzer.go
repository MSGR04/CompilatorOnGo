package Semantic

import (
	"CompilatorOnGo/Core/Parser/Ast"
	"fmt"
)

// SemanticAnalyzer выполняет семантический анализ AST и накапливает ошибки.
// (аналог SemanticAnalyzer из C# версии компилятора)
type SemanticAnalyzer struct {
	environment *SemanticEnvironment
	errors      []string
	warnings    []string
}

// NewSemanticAnalyzer создает новый анализатор.
func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		environment: NewSemanticEnvironment(nil),
		errors:      []string{},
		warnings:    []string{},
	}
}

// Analyze выполняет семантический анализ набора инструкций.
func (a *SemanticAnalyzer) Analyze(statements []Ast.Statement) {
	for _, statement := range statements {
		a.VisitStatement(statement)
	}

	// Проверяем на unused-переменные в глобальной (верхнеуровневой) области.
	for _, name := range a.environment.CollectUnused() {
		a.warnings = append(a.warnings, fmt.Sprintf("Variable '%s' is declared but never used.", name))
	}
}

// VisitStatement обходит инструкцию и проверяет её.
func (a *SemanticAnalyzer) VisitStatement(statement Ast.Statement) {
	switch s := statement.(type) {
	case *Ast.VarStatement:
		if s.Initializer != nil {
			a.VisitExpression(s.Initializer)
		}
		if !a.environment.DefineVariable(s.Name, s.Initializer != nil) {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is already defined.", s.Name))
		}

	case *Ast.PrintStatement:
		a.VisitExpression(s.Expr)

	case *Ast.ExpressionStatement:
		a.VisitExpression(s.Expr)

	case *Ast.BlockStatement:
		previous := a.environment
		a.environment = NewSemanticEnvironment(previous)
		for _, inner := range s.Statements {
			a.VisitStatement(inner)
		}

		// После обработки блока проверяем, что все объявленные в нём переменные были использованы.
		for _, name := range a.environment.CollectUnused() {
			a.warnings = append(a.warnings, fmt.Sprintf("Variable '%s' is declared but never used.", name))
		}

		a.environment = previous

	case *Ast.IfStatement:
		a.VisitExpression(s.Condition)
		a.VisitStatement(s.ThenBranch)
		if s.ElseBranch != nil {
			a.VisitStatement(s.ElseBranch)
		}

	case *Ast.WhileStatement:
		a.VisitExpression(s.Condition)
		a.VisitStatement(s.Body)

	default:
		a.errors = append(a.errors, fmt.Sprintf("Unsupported statement type: %T", statement))
	}
}

// VisitExpression обходит выражение и проверяет его на семантические ошибки.
func (a *SemanticAnalyzer) VisitExpression(expression Ast.Expression) {
	switch e := expression.(type) {
	case *Ast.NumberExpression, *Ast.StringExpression:
		// Ничего не делаем

	case *Ast.VariableExpression:
		defined, initialized := a.environment.UseVariable(e.Name)
		if !defined {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is not defined.", e.Name))
		} else if !initialized {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is used before initialization.", e.Name))
		}

	case *Ast.AssignExpression:
		a.VisitExpression(e.Value)
		if !a.environment.AssignVariable(e.Name) {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is not defined.", e.Name))
		}

	case *Ast.BinaryExpression:
		a.VisitExpression(e.Left)
		a.VisitExpression(e.Right)

	case *Ast.UnaryExpression:
		a.VisitExpression(e.Right)

	default:
		a.errors = append(a.errors, fmt.Sprintf("Unsupported expression type: %T", expression))
	}
}

// Errors возвращает список найденных семантических ошибок.
func (a *SemanticAnalyzer) Errors() []string {
	return a.errors
}

// Warnings возвращает список семантических предупреждений.
func (a *SemanticAnalyzer) Warnings() []string {
	return a.warnings
}
