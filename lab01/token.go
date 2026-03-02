package lab01

import "fmt"

// Token описывает токен: его тип (TokenType) и значение, и позицию в исходном тексте.
type Token struct {
	Type     TokenType
	Value    string
	Position int
}

func NewToken(t TokenType, value string, position int) Token {
	return Token{
		Type:     t,
		Value:    value,
		Position: position,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Token(Type: %v, Value: '%s') at %d", t.Type, t.Value, t.Position)
}
