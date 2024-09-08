package parser

import (
	"fmt"

	"interpreter/internal/expression"
	"interpreter/internal/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]expression.Stmt, error) {
	statements := []expression.Stmt{}

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			fmt.Printf("Error occured: %s", err)
		} else {
			statements = append(statements, stmt)
		}
	}
	if !p.isAtEnd() {
		return nil, fmt.Errorf("unexpected token: %v", p.peek())
	}

	return statements, nil
}
func (p *Parser) assignment() (expression.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.COMMA) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	if p.match(token.QUESTION_MARK) {
		trueExpr, err := p.Expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(token.COLON, "Expect ':' in ternary expression"); err != nil {
			return nil, err
		}

		falseExpr, err := p.Expression()
		if err != nil {
			return nil, err
		}

		expr = expression.NewTernary(expr, trueExpr, falseExpr)
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*expression.Variable); ok {
			name := v.Name

			return expression.NewAssign(name, value), nil
		}

		return nil, fmt.Errorf("invalid assignment target at %s", equals)
	}

	return expr, nil
}
func (p *Parser) Expression() (expression.Expr, error) {
	return p.assignment()
}

func (p *Parser) equality() (expression.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (expression.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (expression.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (expression.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (expression.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expression.NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (expression.Expr, error) {
	if p.match(token.FALSE) {
		return expression.NewLiteral(false), nil
	}
	if p.match(token.TRUE) {
		return expression.NewLiteral(true), nil
	}
	if p.match(token.NIL) {
		return expression.NewLiteral(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return expression.NewLiteral(p.previous().Literal), nil
	}

	if p.match(token.IDENTIFIER) {
		return expression.NewVariable(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return expression.NewGrouping(expr), nil
	}

	return nil, fmt.Errorf("unexpected token: %v", p.peek())
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return token.Token{}, fmt.Errorf("%s at line %d", message, p.peek().Line)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
