package Lexer

import (
	"fmt"
	"unicode"
)

// Предполагается, что TokenType объявлен где-то рядом (или импортируется).
// Предполагается, что Token объявлен где-то рядом.
// Ниже — ожидаемая форма Token:
//
// type Token struct {
//     Type   TokenType
//     Value  string
//     Pos    int
//     Line   int
//     Column int
// }

type Lexer struct {
	input  []rune
	pos    int
	line   int
	column int
}

var keywords = map[string]TokenType{
	"var":   VAR,
	"print": PRINT,
	"if":    IF,
	"else":  ELSE,
	"while": WHILE,
}

var operators = map[string]TokenType{
	"==": EQEQ,
	"!=": NEQ,
	"<=": LTEQ,
	">=": GTEQ,
	"&&": AND,
	"||": OR,
	"+":  PLUS,
	"-":  MINUS,
	"*":  STAR,
	"/":  SLASH,
	"=":  EQ,
	"<":  LT,
	">":  GT,
	"!":  EXCL,
	"(":  LPAREN,
	")":  RPAREN,
	"{":  LBRACE,
	"}":  RBRACE,
	";":  SEMICOLON,
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  []rune(input),
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0)

	for l.pos < len(l.input) {
		current := l.peek()

		if unicode.IsSpace(current) {
			l.next()
			continue
		}

		if unicode.IsDigit(current) {
			tokens = append(tokens, l.readNumber())
			continue
		}

		if unicode.IsLetter(current) {
			tokens = append(tokens, l.readWord())
			continue
		}

		tok, err := l.readOperatorOrPunctuation()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
	}

	// EOF токен: Value "\0" как в C#
	tokens = append(tokens, NewToken(EOF, "\x00", l.pos, l.line, l.column))
	return tokens, nil
}

func (l *Lexer) readNumber() Token {
	startPos := l.pos
	startLine := l.line
	startCol := l.column

	for unicode.IsDigit(l.peek()) {
		l.next()
	}

	text := string(l.input[startPos:l.pos])
	return NewToken(NUMBER, text, startPos, startLine, startCol)
}

func (l *Lexer) readWord() Token {
	startPos := l.pos
	startLine := l.line
	startCol := l.column

	for {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			l.next()
			continue
		}
		break
	}

	text := string(l.input[startPos:l.pos])
	tt, ok := keywords[text]
	if !ok {
		tt = ID
	}

	return NewToken(tt, text, startPos, startLine, startCol)
}

func (l *Lexer) readOperatorOrPunctuation() (Token, error) {
	startPos := l.pos
	startLine := l.line
	startCol := l.column

	// Пробуем 2-символьный оператор
	if l.pos+1 < len(l.input) {
		two := string(l.input[l.pos : l.pos+2])
		if tt, ok := operators[two]; ok {
			l.next()
			l.next()
			return NewToken(tt, two, startPos, startLine, startCol), nil
		}
	}

	// Пробуем 1-символьный
	one := string(l.input[l.pos])
	if tt, ok := operators[one]; ok {
		l.next()
		return NewToken(tt, one, startPos, startLine, startCol), nil
	}

	bad := l.peek()
	return Token{}, fmt.Errorf("[Lexer Error] Unexpected character '%c' at Line %d, Column %d", bad, startLine, startCol)
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	ch := l.input[l.pos]
	l.pos++

	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return ch
}
