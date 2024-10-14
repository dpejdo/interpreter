package expression
import Token "interpreter/internal/token"

type ExprVisitor interface {
    VisitAssignExpr(expr *Assign) interface{}
    VisitBinaryExpr(expr *Binary) interface{}
    VisitTernaryExpr(expr *Ternary) interface{}
    VisitGroupingExpr(expr *Grouping) interface{}
    VisitLiteralExpr(expr *Literal) interface{}
    VisitLogicalExpr(expr *Logical) interface{}
    VisitUnaryExpr(expr *Unary) interface{}
    VisitVariableExpr(expr *Variable) interface{}
}

type Expr interface{
	Accept(visitor ExprVisitor) interface{}
}

type Assign struct {
    Name Token.Token
    Value Expr
}

func NewAssign(Name Token.Token, Value Expr) *Assign {
    return &Assign{
        Name: Name,
        Value: Value,
    }
}

func (e *Assign) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitAssignExpr(e)
}

type Binary struct {
    Left Expr
    Operator Token.Token
    Right Expr
}

func NewBinary(Left Expr, Operator Token.Token, Right Expr) *Binary {
    return &Binary{
        Left: Left,
        Operator: Operator,
        Right: Right,
    }
}

func (e *Binary) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitBinaryExpr(e)
}

type Ternary struct {
    Condition Expr
    TrueExpression Expr
    FalseExpression Expr
}

func NewTernary(Condition Expr, TrueExpression Expr, FalseExpression Expr) *Ternary {
    return &Ternary{
        Condition: Condition,
        TrueExpression: TrueExpression,
        FalseExpression: FalseExpression,
    }
}

func (e *Ternary) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitTernaryExpr(e)
}

type Grouping struct {
    Expr Expr
}

func NewGrouping(Expr Expr) *Grouping {
    return &Grouping{
        Expr: Expr,
    }
}

func (e *Grouping) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitGroupingExpr(e)
}

type Literal struct {
    Value interface{}
}

func NewLiteral(Value interface{}) *Literal {
    return &Literal{
        Value: Value,
    }
}

func (e *Literal) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitLiteralExpr(e)
}

type Logical struct {
    Left Expr
    Operator Token.Token
    Right Expr
}

func NewLogical(Left Expr, Operator Token.Token, Right Expr) *Logical {
    return &Logical{
        Left: Left,
        Operator: Operator,
        Right: Right,
    }
}

func (e *Logical) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitLogicalExpr(e)
}

type Unary struct {
    Operator Token.Token
    Right Expr
}

func NewUnary(Operator Token.Token, Right Expr) *Unary {
    return &Unary{
        Operator: Operator,
        Right: Right,
    }
}

func (e *Unary) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitUnaryExpr(e)
}

type Variable struct {
    Name Token.Token
}

func NewVariable(Name Token.Token) *Variable {
    return &Variable{
        Name: Name,
    }
}

func (e *Variable) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitVariableExpr(e)
}

