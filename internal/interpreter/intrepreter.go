package interpreter

import (
	"fmt"
	environment "interpreter/internal/environment"
	"interpreter/internal/expression"
	"interpreter/internal/token"
)

type Interpreter struct {
	hadError    bool
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: environment.NewEnvironment(nil)}
}

func (i *Interpreter) Interpret(statements []expression.Stmt) interface{} {
	i.hadError = false
	for _, v := range statements {
		i.execute(v)
	}

	/* if i.hadError {
		return nil
	}
	return value */
	return nil
}

func (i *Interpreter) Stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if num, ok := obj.(float64); ok {
		text := fmt.Sprintf("%f", num)
		if text[len(text)-2:] == ".0" {
			return text[:len(text)-2]
		}
		return text
	}

	return fmt.Sprintf("%v", obj)
}

func (i *Interpreter) VisitTernaryExpr(expression *expression.Ternary) interface{} {
	panic("unable to add atm")
	/* 	condition := i.evaluate(expression.Condition)
	   	if i.hadError {
	   		return nil
	   	} */
	/*
		if i.isTruthy(condition) {
			return i.evaluate(expression.TrueExpr)
		}
		return i.evaluate(expression.FalseExpr) */
}

func (i *Interpreter) VisitLiteralExpr(expr *expression.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *expression.Unary) interface{} {
	right := i.evaluate(expr.Right)
	if i.hadError {
		return nil
	}

	switch expr.Operator.Type {
	case token.MINUS:
		if i.checkNumberOperand(expr.Operator, right) {
			return -right.(float64)
		}
	case token.BANG:
		return !i.isTruthy(right)
	}

	i.runtimeError(expr.Operator, "Invalid unary operator")
	return nil
}

func (i *Interpreter) isTruthy(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func (i *Interpreter) VisitGroupingExpr(expr *expression.Grouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) evaluate(expr expression.Expr) interface{} {
	return expr.Accept(i)
}
func (i *Interpreter) execute(stmt expression.Stmt) interface{} {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr *expression.Binary) interface{} {
	left := i.evaluate(expr.Left)
	if i.hadError {
		return nil
	}
	right := i.evaluate(expr.Right)
	if i.hadError {
		return nil
	}

	switch expr.Operator.Type {
	case token.MINUS:
		return i.binaryNumberOp(expr.Operator, left, right, func(a, b float64) float64 { return a - b })
	case token.SLASH:
		return i.binaryNumberOp(expr.Operator, left, right, func(a, b float64) float64 { return a / b })
	case token.STAR:
		return i.binaryNumberOp(expr.Operator, left, right, func(a, b float64) float64 { return a * b })
	case token.PLUS:
		return i.binaryNumberOp(expr.Operator, left, right, func(a, b float64) float64 { return a + b })
	case token.GREATER:
		return i.binaryComparisonOp(expr.Operator, left, right, func(a, b float64) bool { return a > b })
	case token.GREATER_EQUAL:
		return i.binaryComparisonOp(expr.Operator, left, right, func(a, b float64) bool { return a >= b })
	case token.LESS:
		return i.binaryComparisonOp(expr.Operator, left, right, func(a, b float64) bool { return a < b })
	case token.LESS_EQUAL:
		return i.binaryComparisonOp(expr.Operator, left, right, func(a, b float64) bool { return a <= b })
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	i.runtimeError(expr.Operator, "Invalid binary operator")
	return nil
}

func (i *Interpreter) binaryNumberOp(operator token.Token, left, right interface{}, op func(float64, float64) float64) interface{} {
	if i.checkNumberOperands(operator, left, right) {
		return op(left.(float64), right.(float64))
	}
	return nil
}

func (i *Interpreter) binaryComparisonOp(operator token.Token, left, right interface{}, op func(float64, float64) bool) interface{} {
	if i.checkNumberOperands(operator, left, right) {
		return op(left.(float64), right.(float64))
	}
	return nil
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

func (i *Interpreter) checkNumberOperand(operator token.Token, expr interface{}) bool {
	_, ok := expr.(float64)
	if !ok {
		i.runtimeError(operator, "Operand must be a number")
		return false
	}
	return true
}

func (i *Interpreter) checkNumberOperands(operator token.Token, left interface{}, right interface{}) bool {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if !leftOk || !rightOk {
		i.runtimeError(operator, "Operands must be numbers")
		return false
	}
	return true
}

func (i *Interpreter) runtimeError(token token.Token, message string) {
	i.hadError = true
	fmt.Printf("Runtime error at line %d: %s\n", token.Line, message)
}

func (i *Interpreter) VisitExpressionStmt(stmt *expression.Expression) interface{} {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *expression.Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(value)
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

func (i *Interpreter) VisitVariableExpr(expr *expression.Variable) interface{} {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *expression.Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, expr.Value)
	return value

}

func (i *Interpreter) VisitBlockStmt(stmt *expression.Block) interface{} {
	return i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
}
func (i *Interpreter) executeBlock(statements []expression.Stmt, env *environment.Environment) interface{} {
	previous := i.environment
	i.environment = env

	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}
