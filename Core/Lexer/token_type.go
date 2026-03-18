package Lexer

type TokenType int

const (
	NUMBER TokenType = iota
	ID
	STRING
	VAR

	PRINT
	IF
	ELSE
	WHILE

	// Operators
	PLUS  // +
	MINUS // -
	STAR  // *
	SLASH // /
	EQ    // =
	EQEQ  // ==
	EXCL  // !
	NEQ   // !=
	LT    // <
	GT    // >
	LTEQ  // <=
	GTEQ  // >=
	AND   // &&
	OR    // ||

	// Grouping & Punctuation
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	SEMICOLON // ;

	EOF
)
