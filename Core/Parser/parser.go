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
		statement, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func (p *Parser) parseDeclaration() (Ast.Statement, error) {
	if p.match(Lexer.FUNC) {
		return p.parseFunctionDeclaration()
	}
	if p.match(Lexer.VAR) {
		return p.parseVarDeclaration()
	}
	return p.parseStatement()
}

func (p *Parser) parseFunctionDeclaration() (Ast.Statement, error) {
	name, err := p.consume(Lexer.ID, "Expected function name.")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.LPAREN, "Expected '(' after function name."); err != nil {
		return nil, err
	}

	params := make([]string, 0)
	if !p.check(Lexer.RPAREN) {
		for {
			param, err := p.consume(Lexer.ID, "Expected parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, param.Value)

			if !p.match(Lexer.COMMA) {
				break
			}
		}
	}

	if _, err := p.consume(Lexer.RPAREN, "Expected ')' after function parameters."); err != nil {
		return nil, err
	}
	if _, err := p.consume(Lexer.LBRACE, "Expected '{' before function body."); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &Ast.FunctionStatement{
		Name:   name.Value,
		Params: params,
		Body:   body,
	}, nil
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
	if p.match(Lexer.RETURN) {
		return p.parseReturnStatement()
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
	name, err := p.consume(Lexer.ID, "Expected variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Ast.Expression
	if p.match(Lexer.EQ) {
		initializer, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Expected ';' after variable declaration."); err != nil {
		return nil, err
	}

	return &Ast.VarStatement{Name: name.Value, Initializer: initializer}, nil
}

func (p *Parser) parseIfStatement() (Ast.Statement, error) {
	if _, err := p.consume(Lexer.LPAREN, "Expected '(' after 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.RPAREN, "Expected ')' after if condition."); err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch Ast.Statement
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
	if _, err := p.consume(Lexer.LPAREN, "Expected '(' after 'while'."); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.RPAREN, "Expected ')' after while condition."); err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &Ast.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) parsePrintStatement() (Ast.Statement, error) {
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Expected ';' after value."); err != nil {
		return nil, err
	}

	return &Ast.PrintStatement{Expr: value}, nil
}

func (p *Parser) parseReturnStatement() (Ast.Statement, error) {
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Expected ';' after return value."); err != nil {
		return nil, err
	}

	return &Ast.ReturnStatement{Value: value}, nil
}

func (p *Parser) parseExpressionStatement() (Ast.Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Lexer.SEMICOLON, "Expected ';' after expression."); err != nil {
		return nil, err
	}

	return &Ast.ExpressionStatement{Expr: expr}, nil
}

func (p *Parser) parseBlock() ([]Ast.Statement, error) {
	statements := make([]Ast.Statement, 0)

	for !p.check(Lexer.RBRACE) && !p.isAtEnd() {
		statement, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	if _, err := p.consume(Lexer.RBRACE, "Expected '}' after block."); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) parseExpression() (Ast.Expression, error) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (Ast.Expression, error) {
	expr, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	if p.match(Lexer.EQ) {
		equals := p.previous()

		value, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*Ast.VariableExpression); ok {
			return &Ast.AssignExpression{Name: variable.Name, Value: value}, nil
		}

		return nil, fmt.Errorf("[Parser Error] Line %d: invalid assignment target", equals.Line)
	}

	return expr, nil
}

func (p *Parser) parseLogicalOr() (Ast.Expression, error) {
	expr, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.OR) {
		operator := p.previous().Type
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseLogicalAnd() (Ast.Expression, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.AND) {
		operator := p.previous().Type
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseEquality() (Ast.Expression, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.EQEQ, Lexer.NEQ) {
		operator := p.previous().Type
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (Ast.Expression, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.LT, Lexer.LTEQ, Lexer.GT, Lexer.GTEQ) {
		operator := p.previous().Type
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Ast.Expression, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.PLUS, Lexer.MINUS) {
		operator := p.previous().Type
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseFactor() (Ast.Expression, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.match(Lexer.STAR, Lexer.SLASH) {
		operator := p.previous().Type
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		expr = &Ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (Ast.Expression, error) {
	if p.match(Lexer.EXCL, Lexer.MINUS) {
		operator := p.previous().Type
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &Ast.UnaryExpression{Operator: operator, Right: right}, nil
	}
	return p.parseCall()
}

func (p *Parser) parseCall() (Ast.Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if !p.match(Lexer.LPAREN) {
			break
		}

		arguments := make([]Ast.Expression, 0)
		if !p.check(Lexer.RPAREN) {
			for {
				argument, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				arguments = append(arguments, argument)

				if !p.match(Lexer.COMMA) {
					break
				}
			}
		}

		if _, err := p.consume(Lexer.RPAREN, "Expected ')' after arguments."); err != nil {
			return nil, err
		}

		expr = &Ast.CallExpression{Callee: expr, Arguments: arguments}
	}

	return expr, nil
}

func (p *Parser) parsePrimary() (Ast.Expression, error) {
	if p.match(Lexer.NUMBER) {
		raw := p.previous().Value
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			token := p.previous()
			return nil, fmt.Errorf("[Parser Error] Line %d, Col %d: invalid number '%s'", token.Line, token.Column, raw)
		}
		return &Ast.NumberExpression{Value: value}, nil
	}

	if p.match(Lexer.STRING) {
		return &Ast.StringExpression{Value: p.previous().Value}, nil
	}

	if p.match(Lexer.BOOLEAN) {
		return &Ast.BooleanExpression{Value: p.previous().Value == "true"}, nil
	}

	if p.match(Lexer.ID) {
		return &Ast.VariableExpression{Name: p.previous().Value}, nil
	}

	if p.match(Lexer.LPAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(Lexer.RPAREN, "Expected ')' after expression."); err != nil {
			return nil, err
		}
		return expr, nil
	}

	token := p.peek()
	return nil, fmt.Errorf("[Parser Error] Line %d, Col %d: expected expression", token.Line, token.Column)
}

func (p *Parser) match(types ...Lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType Lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
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

func (p *Parser) consume(tokenType Lexer.TokenType, message string) (Lexer.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	token := p.peek()
	return Lexer.Token{}, fmt.Errorf("[Parser Error] Line %d, Col %d: %s", token.Line, token.Column, message)
}
