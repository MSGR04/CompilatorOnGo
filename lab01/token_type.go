package lab01

type TokenType int

const (
	// Literals / identifiers / keywords
	NUMBER TokenType = iota
	ID
	STRING
	VAR

	PRINT
	IF
	ELSE
	WHILE

	// Operators
	PLUS
	MINUS
	STAR
	SLASH
	EQ
	EQEQ
	EXCL
	NEQ
	LT
	GT
	LTEQ
	GTEQ
	AND
	OR

	// Grouping & punctuation
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	SEMICOLON

	EOF
)
