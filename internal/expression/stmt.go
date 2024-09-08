package expression

import Token "interpreter/internal/token"

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitBlockStmt(stmt *Block) interface{}
}

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type Expression struct {
	Expr Expr
}

func NewExpression(expr Expr) *Expression {
	return &Expression{
		Expr: expr,
	}
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func NewPrint(Expression Expr) *Print {
	return &Print{
		Expression: Expression,
	}
}

func (e *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(e)
}

type Var struct {
	Name        Token.Token
	Initializer Expr
}

func NewVar(Name Token.Token, Initializer Expr) *Var {
	return &Var{
		Name:        Name,
		Initializer: Initializer,
	}
}

func (e *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(e)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

func (e *Block) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(e)
}
