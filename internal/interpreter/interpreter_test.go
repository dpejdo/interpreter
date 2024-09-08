package interpreter

import (
	"interpreter/internal/expression"
	"interpreter/internal/token"
	"testing"
)

func TestInterpreter_Interpret(t *testing.T) {
	tests := []struct {
		name    string
		expr    expression.Expr
		want    interface{}
		wantErr bool
	}{
		{
			name: "Simple addition",
			expr: &expression.Binary{
				Left:     &expression.Literal{Value: float64(1)},
				Operator: token.Token{Type: token.PLUS},
				Right:    &expression.Literal{Value: float64(2)},
			},
			want: float64(3),
		},
		{
			name: "Unary negation",
			expr: &expression.Unary{
				Operator: token.Token{Type: token.MINUS},
				Right:    &expression.Literal{Value: float64(5)},
			},
			want: float64(-5),
		},
		{
			name: "Grouping",
			expr: &expression.Grouping{
				Expr: &expression.Literal{Value: float64(42)},
			},
			want: float64(42),
		},
		{
			name: "Comparison",
			expr: &expression.Binary{
				Left:     &expression.Literal{Value: float64(10)},
				Operator: token.Token{Type: token.GREATER},
				Right:    &expression.Literal{Value: float64(5)},
			},
			want: true,
		},
		{
			name: "Equality",
			expr: &expression.Binary{
				Left:     &expression.Literal{Value: float64(5)},
				Operator: token.Token{Type: token.EQUAL_EQUAL},
				Right:    &expression.Literal{Value: float64(5)},
			},
			want: true,
		},
		/* {
			name: "Ternary",
			expr: &expression.Ternary{
				Condition: &expression.Literal{Value: true},
				TrueExpr:  &expression.Literal{Value: float64(1)},
				FalseExpr: &expression.Literal{Value: float64(2)},
			},
			want: float64(1),
		}, */
		{
			name: "Invalid binary operation",
			expr: &expression.Binary{
				Left:     &expression.Literal{Value: "not a number"},
				Operator: token.Token{Type: token.PLUS},
				Right:    &expression.Literal{Value: float64(5)},
			},
			wantErr: true,
		},
		{
			name: "parenthesis binary operation",
			expr: &expression.Binary{
				Left:     expression.NewGrouping(expression.NewBinary(&expression.Literal{Value: float64(5)}, token.Token{Type: token.PLUS}, &expression.Literal{Value: float64(1)})),
				Operator: token.Token{Type: token.STAR},
				Right:    &expression.Literal{Value: float64(2)},
			},
			want: float64(12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInterpreter()
			got := i.Interpret(tt.expr)
			if (i.hadError) != tt.wantErr {
				t.Errorf("Interpreter.Interpret() error = %v, wantErr %v", i.hadError, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Interpreter.Interpret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterpreter_Stringify(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
		want string
	}{
		{"Nil", nil, "nil"},
		{"Integer", float64(42), "42"},
		{"Float", float64(3.14), "3.140000"},
		{"String", "hello", "hello"},
		{"Boolean", true, "true"},
	}

	i := NewInterpreter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := i.Stringify(tt.obj); got != tt.want {
				t.Errorf("Interpreter.Stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}
