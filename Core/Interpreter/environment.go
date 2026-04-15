package Interpreter

import "fmt"

type Environment struct {
	parent *Environment
	values map[string]Value
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		values: make(map[string]Value),
	}
}

func (e *Environment) Define(name string, value Value) {
	e.values[name] = value
}

func (e *Environment) Get(name string) (Value, error) {
	if value, ok := e.values[name]; ok {
		return value, nil
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return Value{}, fmt.Errorf("runtime error: variable '%s' is not defined", name)
}

func (e *Environment) Assign(name string, value Value) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}

	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	return fmt.Errorf("runtime error: variable '%s' is not defined", name)
}
