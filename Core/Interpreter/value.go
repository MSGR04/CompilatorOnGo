package Interpreter

import (
	"fmt"
	"math"

	"CompilatorOnGo/Core/Semantic"
)

type Value struct {
	Type DataType
	Data any
}

type DataType = Semantic.ValueType

func NewNumberValue(value float64) Value {
	return Value{Type: Semantic.NumberType, Data: value}
}

func NewStringValue(value string) Value {
	return Value{Type: Semantic.StringType, Data: value}
}

func NewBooleanValue(value bool) Value {
	return Value{Type: Semantic.BooleanType, Data: value}
}

func NewFunctionValue(function Callable) Value {
	return Value{Type: Semantic.FunctionType, Data: function}
}

func (v Value) Number() (float64, bool) {
	number, ok := v.Data.(float64)
	return number, ok && v.Type == Semantic.NumberType
}

func (v Value) StringValue() (string, bool) {
	text, ok := v.Data.(string)
	return text, ok && v.Type == Semantic.StringType
}

func (v Value) Boolean() (bool, bool) {
	boolean, ok := v.Data.(bool)
	return boolean, ok && v.Type == Semantic.BooleanType
}

func (v Value) Callable() (Callable, bool) {
	callable, ok := v.Data.(Callable)
	return callable, ok && v.Type == Semantic.FunctionType
}

func (v Value) String() string {
	switch v.Type {
	case Semantic.NumberType:
		number, _ := v.Number()
		if math.Trunc(number) == number {
			return fmt.Sprintf("%.0f", number)
		}
		return fmt.Sprintf("%g", number)
	case Semantic.StringType:
		text, _ := v.StringValue()
		return text
	case Semantic.BooleanType:
		boolean, _ := v.Boolean()
		if boolean {
			return "true"
		}
		return "false"
	case Semantic.FunctionType:
		callable, _ := v.Callable()
		return callable.String()
	default:
		return "<invalid>"
	}
}
