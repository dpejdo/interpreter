package expression

import (
	"testing"

	"interpreter/internal/assert"
	Token "interpreter/internal/token"
)

func TestAstPrinter(t *testing.T) {
	tests := []struct {
		name       string
		expression Expr
		expected   string
	}{
		{
			name: "Binary Expression",
			expression: NewBinary(
				NewUnary(
					Token.NewToken(Token.MINUS, "-", nil, 1),
					NewLiteral(123),
				),
				Token.NewToken(Token.STAR, "*", nil, 1),
				NewGrouping(
					NewLiteral(45.67),
				),
			),
			expected: "(* (- 123) (group 45.67))",
		},
		{
			name:       "Grouping Expression",
			expression: NewGrouping(NewLiteral(42)),
			expected:   "(group 42)",
		},
		{
			name:       "Literal Expression",
			expression: NewLiteral("hello"),
			expected:   "hello",
		},
		{
			name:       "Unary Expression",
			expression: NewUnary(Token.NewToken(Token.MINUS, "-", nil, 1), NewLiteral(5)),
			expected:   "(- 5)",
		},
	}

	printer := &AstPrinter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := printer.Print(tt.expression)
			assert.Equal(t, result, tt.expected)
		})
	}
}
