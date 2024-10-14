package expression
import Token "interpreter/internal/token"

type StmtVisitor interface {
    VisitExpressionStmt(stmt *Expression) interface{}
    VisitPrintStmt(stmt *Print) interface{}
    VisitVarStmt(stmt *Var) interface{}
    VisitWhileStmt(stmt *While) interface{}
    VisitBlockStmt(stmt *Block) interface{}
    VisitIfStmt(stmt *If) interface{}
}

type Stmt interface{
	Accept(visitor StmtVisitor) interface{}
}

type Expression struct {
    Expr Expr
}

func NewExpression(Expr Expr) *Expression {
    return &Expression{
        Expr: Expr,
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
    Name Token.Token
    Initializer Expr
}

func NewVar(Name Token.Token, Initializer Expr) *Var {
    return &Var{
        Name: Name,
        Initializer: Initializer,
    }
}

func (e *Var) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitVarStmt(e)
}

type While struct {
    Condition Expr
    Body Stmt
}

func NewWhile(Condition Expr, Body Stmt) *While {
    return &While{
        Condition: Condition,
        Body: Body,
    }
}

func (e *While) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitWhileStmt(e)
}

type Block struct {
    Statements []Stmt
}

func NewBlock(Statements []Stmt) *Block {
    return &Block{
        Statements: Statements,
    }
}

func (e *Block) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitBlockStmt(e)
}

type If struct {
    Condition Expr
    ThenBranch Stmt
    ElseBranch Stmt
}

func NewIf(Condition Expr, ThenBranch Stmt, ElseBranch Stmt) *If {
    return &If{
        Condition: Condition,
        ThenBranch: ThenBranch,
        ElseBranch: ElseBranch,
    }
}

func (e *If) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitIfStmt(e)
}

