package Lexer

import "fmt"

type Token struct {
	Type     TokenType
	Value    string
	Position int
	Line     int
	Column   int
}

func NewToken(t TokenType, value string, position, line, column int) Token {
	return Token{
		Type:     t,
		Value:    value,
		Position: position,
		Line:     line,
		Column:   column,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("[%d:%d] Token(%v, '%s')", t.Line, t.Column, t.Type, t.Value)
}
