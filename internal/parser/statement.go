package parser

import (
	"interpreter/internal/expression"
	"interpreter/internal/token"
)

func (p *Parser) Statement() (expression.Stmt, error) {
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.LEFT_BRACE) {
		if val, err := p.block(); err == nil {

			return expression.NewBlock(val), nil
		} else {
			return nil, err
		}
	}

	return p.expressionStatement()
}

func (p *Parser) Declaration() (expression.Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	stmt, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return stmt, nil

}
func (p *Parser) varDeclaration() (expression.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer expression.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration\n")
	if err != nil {
		return nil, err
	}
	return expression.NewVar(name, initializer), nil
}
func (p *Parser) whileStatement() (expression.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.Statement()
	if err != nil {
		return nil, err
	}
	return expression.NewWhile(condition, body), nil

}
func (p *Parser) forStatement() (expression.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	var initializer expression.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		value, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = value
	} else {
		value, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = value
	}
	var condition expression.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after loop condition.")
	if err != nil {
		return nil, err
	}
	var increment expression.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after for loop increment.")
	if err != nil {
		return nil, err
	}

	body, err := p.Statement()

	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = expression.NewBlock([]expression.Stmt{body, expression.NewExpression(increment)})
	}

	if condition == nil {
		condition = expression.NewLiteral(true)
	}
	body = expression.NewWhile(condition, body)

	if initializer != nil {
		body = expression.NewBlock([]expression.Stmt{initializer, body})
	}
	return body, nil
}

func (p *Parser) printStatement() (expression.Stmt, error) {
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return expression.NewPrint(value), nil
}

func (p *Parser) expressionStatement() (expression.Stmt, error) {
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after value.")

	return expression.NewExpression(value), nil
}

func (p *Parser) block() ([]expression.Stmt, error) {
	var stmts []expression.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {

		stmt, err := p.Declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	p.consume(token.RIGHT_BRACE, " Expected '}' after block.\n")

	return stmts, nil
}

func (p *Parser) ifStatement() (expression.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition")

	thenBranch, err := p.Statement()
	if err != nil {
		return nil, err
	}
	var elseBranch expression.Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.Statement()
		if err != nil {
			return nil, err
		}
	}

	return expression.NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) or() (expression.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expression.NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (expression.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expression.NewLogical(expr, operator, right)
	}

	return expr, nil
}
