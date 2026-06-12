package Semantic

import (
	"fmt"

	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser/Ast"
)

type SemanticAnalyzer struct {
	environment     *SemanticEnvironment
	errors          []string
	warnings        []string
	currentFunction *functionInfo
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		environment: NewSemanticEnvironment(nil),
		errors:      []string{},
		warnings:    []string{},
	}
}

func (a *SemanticAnalyzer) Analyze(statements []Ast.Statement) {
	for _, statement := range statements {
		a.VisitStatement(statement)
	}

	for _, name := range a.environment.CollectUnused() {
		a.warnings = append(a.warnings, fmt.Sprintf("'%s' is declared but never used.", name))
	}
}

func (a *SemanticAnalyzer) VisitStatement(statement Ast.Statement) {
	switch s := statement.(type) {
	case *Ast.VarStatement:
		initializerType := UnknownType
		initialized := false

		if s.Initializer == nil {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' must be initialized to infer its static type.", s.Name))
		} else {
			initializerType = a.VisitExpression(s.Initializer)
			initialized = true
		}

		if !a.environment.DefineVariable(s.Name, initializerType, initialized) {
			a.errors = append(a.errors, fmt.Sprintf("Name '%s' is already defined in this scope.", s.Name))
		}

	case *Ast.FunctionStatement:
		function, ok := a.environment.DefineFunction(s.Name, len(s.Params))
		if !ok {
			a.errors = append(a.errors, fmt.Sprintf("Name '%s' is already defined in this scope.", s.Name))
			return
		}

		previousEnvironment := a.environment
		previousFunction := a.currentFunction

		a.environment = NewSemanticEnvironment(previousEnvironment)
		a.currentFunction = function

		seenParams := make(map[string]struct{}, len(s.Params))
		for _, param := range s.Params {
			if _, exists := seenParams[param]; exists {
				a.errors = append(a.errors, fmt.Sprintf("Function '%s' has duplicate parameter '%s'.", s.Name, param))
				continue
			}
			seenParams[param] = struct{}{}

			if !a.environment.DefineVariable(param, UnknownType, true) {
				a.errors = append(a.errors, fmt.Sprintf("Parameter '%s' conflicts with another name in function '%s'.", param, s.Name))
			}
		}

		for _, inner := range s.Body {
			a.VisitStatement(inner)
		}

		for _, name := range a.environment.CollectUnused() {
			a.warnings = append(a.warnings, fmt.Sprintf("'%s' is declared but never used.", name))
		}

		a.environment = previousEnvironment
		a.currentFunction = previousFunction

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

		for _, name := range a.environment.CollectUnused() {
			a.warnings = append(a.warnings, fmt.Sprintf("'%s' is declared but never used.", name))
		}

		a.environment = previous

	case *Ast.IfStatement:
		conditionType := a.VisitExpression(s.Condition)
		if conditionType != UnknownType && conditionType != BooleanType {
			a.errors = append(a.errors, fmt.Sprintf("Condition in if statement must be Boolean, got %s.", conditionType))
		}

		a.VisitStatement(s.ThenBranch)
		if s.ElseBranch != nil {
			a.VisitStatement(s.ElseBranch)
		}

	case *Ast.WhileStatement:
		conditionType := a.VisitExpression(s.Condition)
		if conditionType != UnknownType && conditionType != BooleanType {
			a.errors = append(a.errors, fmt.Sprintf("Condition in while statement must be Boolean, got %s.", conditionType))
		}

		a.VisitStatement(s.Body)

	case *Ast.ReturnStatement:
		if a.currentFunction == nil {
			a.errors = append(a.errors, "Return statement is only allowed inside a function.")
			return
		}

		returnType := a.VisitExpression(s.Value)
		if a.currentFunction.returnType == UnknownType {
			a.currentFunction.returnType = returnType
		} else if returnType != UnknownType && a.currentFunction.returnType != returnType {
			a.errors = append(a.errors, fmt.Sprintf("Function returns inconsistent types: %s and %s.", a.currentFunction.returnType, returnType))
		}

	default:
		a.errors = append(a.errors, fmt.Sprintf("Unsupported statement type: %T", statement))
	}
}

func (a *SemanticAnalyzer) VisitExpression(expression Ast.Expression) ValueType {
	switch e := expression.(type) {
	case *Ast.NumberExpression:
		return NumberType

	case *Ast.StringExpression:
		return StringType

	case *Ast.BooleanExpression:
		return BooleanType

	case *Ast.VariableExpression:
		kind, variable, function := a.environment.ResolveName(e.Name)
		switch kind {
		case symbolVariable:
			variable.used = true
			if !variable.initialized {
				a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is used before initialization.", e.Name))
				return UnknownType
			}
			return variable.valueType
		case symbolFunction:
			function.used = true
			return FunctionType
		default:
			a.errors = append(a.errors, fmt.Sprintf("Name '%s' is not defined.", e.Name))
			return UnknownType
		}

	case *Ast.AssignExpression:
		valueType := a.VisitExpression(e.Value)
		kind, info, _ := a.environment.ResolveName(e.Name)
		if kind == symbolMissing {
			a.errors = append(a.errors, fmt.Sprintf("Variable '%s' is not defined.", e.Name))
			return UnknownType
		}
		if kind == symbolFunction {
			a.errors = append(a.errors, fmt.Sprintf("Cannot assign to function '%s'.", e.Name))
			return UnknownType
		}

		if info.valueType == UnknownType {
			info.valueType = valueType
		} else if valueType != UnknownType && info.valueType != valueType {
			a.errors = append(a.errors, fmt.Sprintf("Cannot assign value of type %s to variable '%s' of type %s.", valueType, e.Name, info.valueType))
			return UnknownType
		}

		info.initialized = true
		return info.valueType

	case *Ast.CallExpression:
		for _, argument := range e.Arguments {
			a.VisitExpression(argument)
		}

		calleeName, ok := getCalledFunctionName(e.Callee)
		if !ok {
			calleeType := a.VisitExpression(e.Callee)
			if calleeType != UnknownType && calleeType != FunctionType {
				a.errors = append(a.errors, fmt.Sprintf("Attempted to call non-function value of type %s.", calleeType))
			}
			return UnknownType
		}

		kind, variable, function := a.environment.ResolveName(calleeName)
		switch kind {
		case symbolVariable:
			variable.used = true
			a.errors = append(a.errors, fmt.Sprintf("Attempted to call variable '%s' of type %s.", calleeName, variable.valueType))
			return UnknownType
		case symbolFunction:
			function.used = true
		default:
			a.errors = append(a.errors, fmt.Sprintf("Function '%s' is not defined.", calleeName))
			return UnknownType
		}

		if len(e.Arguments) != function.paramCount {
			a.errors = append(a.errors, fmt.Sprintf("Function '%s' expects %d arguments, got %d.", calleeName, function.paramCount, len(e.Arguments)))
		}

		return function.returnType

	case *Ast.BinaryExpression:
		leftType := a.VisitExpression(e.Left)
		rightType := a.VisitExpression(e.Right)

		switch e.Operator {
		case Lexer.PLUS:
			if leftType == NumberType && rightType == NumberType {
				return NumberType
			}
			if leftType == StringType && rightType == StringType {
				return StringType
			}
			if leftType != UnknownType && rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Operator '+' expects Number+Number or String+String, got %s and %s.", leftType, rightType))
			}
			return UnknownType

		case Lexer.MINUS, Lexer.STAR, Lexer.SLASH:
			if leftType == NumberType && rightType == NumberType {
				return NumberType
			}
			if leftType != UnknownType && rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Arithmetic operator expects Number operands, got %s and %s.", leftType, rightType))
			}
			return UnknownType

		case Lexer.GT, Lexer.GTEQ, Lexer.LT, Lexer.LTEQ:
			if leftType == NumberType && rightType == NumberType {
				return BooleanType
			}
			if leftType != UnknownType && rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Comparison operator expects Number operands, got %s and %s.", leftType, rightType))
			}
			return UnknownType

		case Lexer.EQEQ, Lexer.NEQ:
			if leftType == UnknownType || rightType == UnknownType {
				return UnknownType
			}
			if leftType != rightType {
				a.errors = append(a.errors, fmt.Sprintf("Equality operator requires operands of the same type, got %s and %s.", leftType, rightType))
				return UnknownType
			}
			return BooleanType

		case Lexer.AND, Lexer.OR:
			if leftType == BooleanType && rightType == BooleanType {
				return BooleanType
			}
			if leftType != UnknownType && rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Logical operator expects Boolean operands, got %s and %s.", leftType, rightType))
			}
			return UnknownType

		default:
			a.errors = append(a.errors, fmt.Sprintf("Unsupported binary operator: %v", e.Operator))
			return UnknownType
		}

	case *Ast.UnaryExpression:
		rightType := a.VisitExpression(e.Right)

		switch e.Operator {
		case Lexer.MINUS:
			if rightType == NumberType {
				return NumberType
			}
			if rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Unary '-' expects Number, got %s.", rightType))
			}
			return UnknownType

		case Lexer.EXCL:
			if rightType == BooleanType {
				return BooleanType
			}
			if rightType != UnknownType {
				a.errors = append(a.errors, fmt.Sprintf("Unary '!' expects Boolean, got %s.", rightType))
			}
			return UnknownType

		default:
			a.errors = append(a.errors, fmt.Sprintf("Unsupported unary operator: %v", e.Operator))
			return UnknownType
		}

	default:
		a.errors = append(a.errors, fmt.Sprintf("Unsupported expression type: %T", expression))
		return UnknownType
	}
}

func getCalledFunctionName(expression Ast.Expression) (string, bool) {
	variable, ok := expression.(*Ast.VariableExpression)
	if !ok {
		return "", false
	}
	return variable.Name, true
}

func (a *SemanticAnalyzer) Errors() []string {
	return a.errors
}

func (a *SemanticAnalyzer) Warnings() []string {
	return a.warnings
}
