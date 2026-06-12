package Interpreter

import (
	"fmt"

	"CompilatorOnGo/Core/Parser/Ast"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []Value) (Value, error)
	String() string
}

type UserFunction struct {
	declaration *Ast.FunctionStatement
	closure     *Environment
}

func NewUserFunction(declaration *Ast.FunctionStatement, closure *Environment) *UserFunction {
	return &UserFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *UserFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *UserFunction) Call(interpreter *Interpreter, arguments []Value) (Value, error) {
	environment := NewEnvironment(f.closure)
	for index, param := range f.declaration.Params {
		environment.Define(param, arguments[index])
	}

	err := interpreter.executeBlock(f.declaration.Body, environment)
	if err == nil {
		return Value{}, fmt.Errorf("runtime error: function '%s' must return a value", f.declaration.Name)
	}

	signal, ok := err.(*returnSignal)
	if !ok {
		return Value{}, err
	}
	return signal.value, nil
}

func (f *UserFunction) String() string {
	return fmt.Sprintf("<function %s>", f.declaration.Name)
}

type returnSignal struct {
	value Value
}

func (r *returnSignal) Error() string {
	return "return"
}
