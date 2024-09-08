package parser

import (
	"interpreter/internal/expression"
	"interpreter/internal/token"
)

func (p *Parser) statement() (expression.Stmt, error) {

	if p.match(token.PRINT) {
		return p.printStatement()
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

func (p *Parser) declaration() (expression.Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	stmt, err := p.statement()
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

		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	p.consume(token.RIGHT_BRACE, " Expected '}' after block.\n")

	return stmts, nil
}
