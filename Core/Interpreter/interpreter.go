package Interpreter

import (
	"fmt"
	"io"
	"os"

	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser/Ast"
)

type Interpreter struct {
	environment *Environment
	output      io.Writer
}

func New(output io.Writer) *Interpreter {
	if output == nil {
		output = os.Stdout
	}

	return &Interpreter{
		environment: NewEnvironment(nil),
		output:      output,
	}
}

func (i *Interpreter) Interpret(statements []Ast.Statement) error {
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(statement Ast.Statement) error {
	switch s := statement.(type) {
	case *Ast.ExpressionStatement:
		_, err := i.evaluate(s.Expr)
		return err

	case *Ast.PrintStatement:
		value, err := i.evaluate(s.Expr)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(i.output, value.String())
		return err

	case *Ast.VarStatement:
		if s.Initializer == nil {
			return fmt.Errorf("runtime error: variable '%s' must be initialized", s.Name)
		}

		value, err := i.evaluate(s.Initializer)
		if err != nil {
			return err
		}

		i.environment.Define(s.Name, value)
		return nil

	case *Ast.BlockStatement:
		return i.executeBlock(s.Statements, NewEnvironment(i.environment))

	case *Ast.IfStatement:
		condition, err := i.evaluate(s.Condition)
		if err != nil {
			return err
		}

		boolean, ok := condition.Boolean()
		if !ok {
			return fmt.Errorf("runtime error: if condition must be Boolean")
		}

		if boolean {
			return i.execute(s.ThenBranch)
		}

		if s.ElseBranch != nil {
			return i.execute(s.ElseBranch)
		}

		return nil

	case *Ast.WhileStatement:
		for {
			condition, err := i.evaluate(s.Condition)
			if err != nil {
				return err
			}

			boolean, ok := condition.Boolean()
			if !ok {
				return fmt.Errorf("runtime error: while condition must be Boolean")
			}

			if !boolean {
				return nil
			}

			if err := i.execute(s.Body); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("runtime error: unsupported statement type %T", statement)
	}
}

func (i *Interpreter) executeBlock(statements []Ast.Statement, environment *Environment) error {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) evaluate(expression Ast.Expression) (Value, error) {
	switch e := expression.(type) {
	case *Ast.NumberExpression:
		return NewNumberValue(e.Value), nil

	case *Ast.StringExpression:
		return NewStringValue(e.Value), nil

	case *Ast.BooleanExpression:
		return NewBooleanValue(e.Value), nil

	case *Ast.VariableExpression:
		return i.environment.Get(e.Name)

	case *Ast.AssignExpression:
		value, err := i.evaluate(e.Value)
		if err != nil {
			return Value{}, err
		}

		if err := i.environment.Assign(e.Name, value); err != nil {
			return Value{}, err
		}

		return value, nil

	case *Ast.UnaryExpression:
		right, err := i.evaluate(e.Right)
		if err != nil {
			return Value{}, err
		}

		switch e.Operator {
		case Lexer.MINUS:
			number, ok := right.Number()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: unary '-' expects Number")
			}
			return NewNumberValue(-number), nil

		case Lexer.EXCL:
			boolean, ok := right.Boolean()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: unary '!' expects Boolean")
			}
			return NewBooleanValue(!boolean), nil

		default:
			return Value{}, fmt.Errorf("runtime error: unsupported unary operator %v", e.Operator)
		}

	case *Ast.BinaryExpression:
		switch e.Operator {
		case Lexer.OR:
			left, err := i.evaluate(e.Left)
			if err != nil {
				return Value{}, err
			}

			leftBoolean, ok := left.Boolean()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: operator '||' expects Boolean operands")
			}

			if leftBoolean {
				return NewBooleanValue(true), nil
			}

			right, err := i.evaluate(e.Right)
			if err != nil {
				return Value{}, err
			}

			rightBoolean, ok := right.Boolean()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: operator '||' expects Boolean operands")
			}

			return NewBooleanValue(rightBoolean), nil

		case Lexer.AND:
			left, err := i.evaluate(e.Left)
			if err != nil {
				return Value{}, err
			}

			leftBoolean, ok := left.Boolean()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: operator '&&' expects Boolean operands")
			}

			if !leftBoolean {
				return NewBooleanValue(false), nil
			}

			right, err := i.evaluate(e.Right)
			if err != nil {
				return Value{}, err
			}

			rightBoolean, ok := right.Boolean()
			if !ok {
				return Value{}, fmt.Errorf("runtime error: operator '&&' expects Boolean operands")
			}

			return NewBooleanValue(rightBoolean), nil
		}

		left, err := i.evaluate(e.Left)
		if err != nil {
			return Value{}, err
		}

		right, err := i.evaluate(e.Right)
		if err != nil {
			return Value{}, err
		}

		switch e.Operator {
		case Lexer.PLUS:
			if leftNumber, ok := left.Number(); ok {
				rightNumber, rightOk := right.Number()
				if !rightOk {
					return Value{}, fmt.Errorf("runtime error: '+' expects Number+Number or String+String")
				}
				return NewNumberValue(leftNumber + rightNumber), nil
			}

			leftString, leftOk := left.StringValue()
			rightString, rightOk := right.StringValue()
			if leftOk && rightOk {
				return NewStringValue(leftString + rightString), nil
			}

			return Value{}, fmt.Errorf("runtime error: '+' expects Number+Number or String+String")

		case Lexer.MINUS:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewNumberValue(a - b), nil
			}, "-")

		case Lexer.STAR:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewNumberValue(a * b), nil
			}, "*")

		case Lexer.SLASH:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				if b == 0 {
					return Value{}, fmt.Errorf("runtime error: division by zero")
				}
				return NewNumberValue(a / b), nil
			}, "/")

		case Lexer.GT:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewBooleanValue(a > b), nil
			}, ">")

		case Lexer.GTEQ:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewBooleanValue(a >= b), nil
			}, ">=")

		case Lexer.LT:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewBooleanValue(a < b), nil
			}, "<")

		case Lexer.LTEQ:
			return executeNumericBinary(left, right, func(a, b float64) (Value, error) {
				return NewBooleanValue(a <= b), nil
			}, "<=")

		case Lexer.EQEQ:
			return NewBooleanValue(valuesEqual(left, right)), nil

		case Lexer.NEQ:
			return NewBooleanValue(!valuesEqual(left, right)), nil

		default:
			return Value{}, fmt.Errorf("runtime error: unsupported binary operator %v", e.Operator)
		}

	default:
		return Value{}, fmt.Errorf("runtime error: unsupported expression type %T", expression)
	}
}

func executeNumericBinary(left Value, right Value, operation func(float64, float64) (Value, error), operator string) (Value, error) {
	leftNumber, leftOk := left.Number()
	rightNumber, rightOk := right.Number()
	if !leftOk || !rightOk {
		return Value{}, fmt.Errorf("runtime error: operator '%s' expects Number operands", operator)
	}

	return operation(leftNumber, rightNumber)
}

func valuesEqual(left Value, right Value) bool {
	if left.Type != right.Type {
		return false
	}

	return left.Data == right.Data
}
