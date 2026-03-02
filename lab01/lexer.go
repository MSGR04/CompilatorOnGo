package lab01

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input  []rune
	length int
	pos    int
}

func NewLexer(input string) *Lexer {
	r := []rune(input)
	return &Lexer{
		input:  r,
		length: len(r),
		pos:    0,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	result := make([]Token, 0)

	for l.pos < l.length {
		current := l.peek()

		if unicode.IsSpace(current) {
			l.next()
			continue
		}

		if unicode.IsDigit(current) {
			l.tokenizeNumber(&result)
			continue
		}

		if unicode.IsLetter(current) {
			l.tokenizeWord(&result)
			continue
		}

		if err := l.tokenizeOperator(&result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (l *Lexer) tokenizeNumber(result *[]Token) {
	start := l.pos

	for unicode.IsDigit(l.peek()) {
		l.next()
	}

	numberStr := string(l.input[start:l.pos])
	*result = append(*result, NewToken(NUMBER, numberStr, start))
}

func (l *Lexer) tokenizeWord(result *[]Token) {
	start := l.pos

	for {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			l.next()
			continue
		}
		break
	}

	word := string(l.input[start:l.pos])

	switch word {
	case "var":
		l.addToken(result, VAR, word, start)
	case "print":
		l.addToken(result, PRINT, word, start)
	case "if":
		l.addToken(result, IF, word, start)
	case "else":
		l.addToken(result, ELSE, word, start)
	case "while":
		l.addToken(result, WHILE, word, start)
	default:
		l.addToken(result, ID, word, start)
	}
}

func (l *Lexer) tokenizeOperator(result *[]Token) error {
	current := l.peek()
	start := l.pos

	switch current {
	case '+':
		l.next()
		l.addToken(result, PLUS, "+", start)
	case '-':
		l.next()
		l.addToken(result, MINUS, "-", start)
	case '*':
		l.next()
		l.addToken(result, STAR, "*", start)
	case '/':
		l.next()
		l.addToken(result, SLASH, "/", start)
	case '=':
		l.next()
		l.addToken(result, EQ, "=", start)
	case ';':
		l.next()
		l.addToken(result, SEMICOLON, ";", start)
	default:
		return fmt.Errorf("unexpected character '%c' at position %d", current, l.pos)
	}

	return nil
}

func (l *Lexer) peek() rune {
	if l.pos >= l.length {
		return 0 // аналог '\0'
	}
	return l.input[l.pos]
}

func (l *Lexer) next() rune {
	if l.pos >= l.length {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) addToken(result *[]Token, t TokenType, value string, start int) {
	*result = append(*result, NewToken(t, value, start))
}
