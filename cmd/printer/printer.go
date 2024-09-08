package main

import (
	"fmt"

	tt "interpreter/internal/expression"
	token "interpreter/internal/token"
)

func main() {
	expression := tt.NewBinary(
		tt.NewLiteral(4),
		token.NewToken(token.STAR, "*", nil, 0),
		tt.NewGrouping(
			tt.NewBinary(
				tt.NewLiteral(2),
				token.NewToken(token.PLUS, "+", nil, 0),
				tt.NewLiteral(3),
			),
		),
	)
	printer := &tt.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
