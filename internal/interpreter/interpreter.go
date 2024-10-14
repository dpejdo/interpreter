// File: interpreter.go

package interpreter

import (
	"fmt"
	"interpreter/internal/environment"
	"interpreter/internal/expression"
	"interpreter/internal/token"
)

type Interpreter struct {
	environment *environment.Environment
}

// VisitTernaryExpr implements expression.ExprVisitor.
func (i *Interpreter) VisitTernaryExpr(expr *expression.Ternary) interface{} {
	panic("unimplemented")
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: environment.NewEnvironment(nil)}
}

func (i *Interpreter) Interpret(statements []expression.Stmt) {
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) VisitExpressionStmt(stmt *expression.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *expression.Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *expression.Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *expression.Block) interface{} {
	i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *expression.If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *expression.While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *expression.Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *expression.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.PLUS:
		return i.add(left, right, expr.Operator)
	case token.MINUS:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) float64 { return a - b })
	case token.STAR:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) float64 { return a * b })
	case token.SLASH:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) float64 { return a / b })
	case token.GREATER:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) bool { return a > b })
	case token.GREATER_EQUAL:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) bool { return a >= b })
	case token.LESS:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) bool { return a < b })
	case token.LESS_EQUAL:
		return i.checkNumberOperands(expr.Operator, left, right, func(a, b float64) bool { return a <= b })
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *expression.Grouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *expression.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *expression.Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr *expression.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		return -i.checkNumberOperand(expr.Operator, right)
	case token.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *expression.Variable) interface{} {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) execute(stmt expression.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr expression.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) executeBlock(statements []expression.Stmt, env *environment.Environment) {
	previous := i.environment
	i.environment = env
	defer func() { i.environment = previous }()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand interface{}) float64 {
	if value, ok := operand.(float64); ok {
		return value
	}
	panic(i.runtimeError(operator, "Operand must be a number."))
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left, right interface{}, op interface{}) interface{} {
	leftVal, leftOk := left.(float64)
	rightVal, rightOk := right.(float64)
	if !leftOk || !rightOk {
		panic(i.runtimeError(operator, "Operands must be numbers."))
	}
	switch fn := op.(type) {
	case func(float64, float64) float64:
		return fn(leftVal, rightVal)
	case func(float64, float64) bool:
		return fn(leftVal, rightVal)
	}
	return nil
}

func (i *Interpreter) add(left, right interface{}, operator token.Token) interface{} {
	if leftNum, leftOk := left.(float64); leftOk {
		if rightNum, rightOk := right.(float64); rightOk {
			return leftNum + rightNum
		}
	}
	if leftStr, leftOk := left.(string); leftOk {
		if rightStr, rightOk := right.(string); rightOk {
			return leftStr + rightStr
		}
	}
	panic(i.runtimeError(operator, "Operands must be two numbers or two strings."))
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if num, ok := object.(float64); ok {
		return fmt.Sprintf("%g", num)
	}
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) runtimeError(token token.Token, message string) RuntimeError {
	return RuntimeError{Token: token, Message: message}
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.Message, e.Token.Line)
}
