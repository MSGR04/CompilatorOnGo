package Semantic

type ValueType int

const (
	UnknownType ValueType = iota
	NumberType
	StringType
	BooleanType
)

func (t ValueType) String() string {
	switch t {
	case NumberType:
		return "Number"
	case StringType:
		return "String"
	case BooleanType:
		return "Boolean"
	default:
		return "Unknown"
	}
}
