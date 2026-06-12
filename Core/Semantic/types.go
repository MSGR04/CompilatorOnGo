package Semantic

type ValueType int

const (
	UnknownType ValueType = iota
	NumberType
	StringType
	BooleanType
	FunctionType
)

func (t ValueType) String() string {
	switch t {
	case NumberType:
		return "Number"
	case StringType:
		return "String"
	case BooleanType:
		return "Boolean"
	case FunctionType:
		return "Function"
	default:
		return "Unknown"
	}
}
