package Parser

import (
	"fmt"
	"strconv"

	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser/Ast"
)

type Parser struct {
	tokens   []Lexer.Token
	position int
}

func New(tokens []Lexer.Token) *Parser {
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) Parse() ([]Ast.Statement, error) {
	statements := make([]Ast.Statement, 0)
	for !p.isAtEnd() {
		st, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, st)
	}
	return statements, nil
}

func (p *Parser) parseDeclaration() (Ast.Statement, error) {
	if p.match(Lexer.VAR) {
		return p.parseVarDeclaration()
	}
	return p.parseStatement()
}

func (p *Parser) parseStatement() (Ast.Statement, error) {
	if p.match(Lexer.IF) {
		return p.parseIfStatement()
	}
	if p.match(Lexer.WHILE) {
		return p.parseWhileStatement()
	}
	if p.match(Lexer.PRINT) {
		return p.parsePrintStatement()
	}
	if p.match(Lexer.LBRACE) {
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		return &Ast.BlockStatement{Statements: block}, nil
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseVarDeclaration() (Ast.Statement, error) {
	name, err := p.consume(Lexer.ID, "Ожидается имя переменной.")
	if err != nil {
		return nil, err
	}

	var initializer Ast.Expression = nil
	if p.match(Lexer.EQ) {
		initializer, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Ожидается ';' после объявления переменной."); err != nil {
		return nil, err
	}

	return &Ast.VarStatement{
		Name:        name.Value,
		Initializer: initializer,
	}, nil
}

func (p *Parser) parseIfStatement() (Ast.Statement, error) {
	if _, err := p.consume(Lexer.LPAREN, "Ожидается '(' после 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.RPAREN, "Ожидается ')' после условия 'if'."); err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch Ast.Statement = nil
	if p.match(Lexer.ELSE) {
		elseBranch, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &Ast.IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) parseWhileStatement() (Ast.Statement, error) {
	if _, err := p.consume(Lexer.LPAREN, "Ожидается '(' после 'while'."); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.RPAREN, "Ожидается ')' после условия 'while'."); err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &Ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parsePrintStatement() (Ast.Statement, error) {
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Ожидается ';' после значения."); err != nil {
		return nil, err
	}

	return &Ast.PrintStatement{Expr: value}, nil
}

func (p *Parser) parseExpressionStatement() (Ast.Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Ожидается ';' после выражения."); err != nil {
		return nil, err
	}

	return &Ast.ExpressionStatement{Expr: expr}, nil
}

func (p *Parser) parseBlock() ([]Ast.Statement, error) {
	stmts := make([]Ast.Statement, 0)

	for !p.check(Lexer.RBRACE) && !p.isAtEnd() {
		st, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, st)
	}

	if _, err := p.consume(Lexer.RBRACE, "Ожидается '}' после блока."); err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) parseExpression() (Ast.Expression, error) {
	return p.parseAssignment()
}

// 1) Присваивание (самый низкий приоритет)
func (p *Parser) parseAssignment() (Ast.Expression, error) {
	expr, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	if p.match(Lexer.EQ) {
		equals := p.previous()

		value, err := p.parseAssignment() // a = b = 5
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*Ast.VariableExpression); ok {
			return &Ast.AssignExpression{Name: v.Name, Value: value}, nil
		}

		return nil, fmt.Errorf("[Parser Error] Line %d: Недопустимая цель для присваивания.", equals.Line)
	}

	return expr, nil
}

// 2) ||
func (p *Parser) parseLogicalOr() (Ast.Expression, error) {
	expr, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.OR) {
		op := p.previous().Type
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

// 3) &&
func (p *Parser) parseLogicalAnd() (Ast.Expression, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.AND) {
		op := p.previous().Type
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

// 4) ==, !=
func (p *Parser) parseEquality() (Ast.Expression, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.EQEQ, Lexer.NEQ) {
		op := p.previous().Type
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

// 5) <, >, <=, >=
func (p *Parser) parseComparison() (Ast.Expression, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.LT, Lexer.LTEQ, Lexer.GT, Lexer.GTEQ) {
		op := p.previous().Type
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

// 6) +, -
func (p *Parser) parseTerm() (Ast.Expression, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.PLUS, Lexer.MINUS) {
		op := p.previous().Type
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

// 7) *, /
func (p *Parser) parseFactor() (Ast.Expression, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.STAR, Lexer.SLASH) {
		op := p.previous().Type
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (Ast.Expression, error) {
	if p.match(Lexer.EXCL, Lexer.MINUS) {
		op := p.previous().Type
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &Ast.UnaryExpression{Operator: op, Right: right}, nil
	}
	return p.parsePrimary()
}

// 9) примитивы
func (p *Parser) parsePrimary() (Ast.Expression, error) {
	if p.match(Lexer.NUMBER) {
		raw := p.previous().Value
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			tok := p.previous()
			return nil, fmt.Errorf("[Parser Error] Line %d, Col %d: Некорректное число '%s'.", tok.Line, tok.Column, raw)
		}
		return &Ast.NumberExpression{Value: val}, nil
	}

	if p.match(Lexer.STRING) {
		return &Ast.StringExpression{Value: p.previous().Value}, nil
	}

	if p.match(Lexer.ID) {
		return &Ast.VariableExpression{Name: p.previous().Value}, nil
	}

	if p.match(Lexer.LPAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(Lexer.RPAREN, "Ожидается ')' после выражения."); err != nil {
			return nil, err
		}
		return expr, nil
	}

	tok := p.peek()
	return nil, fmt.Errorf("[Parser Error] Line %d, Col %d: Ожидается выражение.", tok.Line, tok.Column)
}

// ----------------- helpers -----------------

func (p *Parser) match(types ...Lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t Lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Lexer.Token {
	if !p.isAtEnd() {
		p.position++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == Lexer.EOF
}

func (p *Parser) peek() Lexer.Token {
	return p.tokens[p.position]
}

func (p *Parser) previous() Lexer.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) consume(t Lexer.TokenType, message string) (Lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	tok := p.peek()
	return Lexer.Token{}, fmt.Errorf("[Parser Error] Line %d, Col %d: %s", tok.Line, tok.Column, message)
}
