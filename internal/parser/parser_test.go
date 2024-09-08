package parser

import (
	"interpreter/internal/expression"
	"interpreter/internal/scanner"
	"interpreter/internal/token"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    expression.Expr
		wantErr bool
		errMsg  string
	}{
		{
			name:  "Simple addition",
			input: "1 + 2",
			want: expression.NewBinary(
				expression.NewLiteral(float64(1)),
				token.Token{Type: token.PLUS, Lexeme: "+"},
				expression.NewLiteral(float64(2)),
			),
			wantErr: false,
		},
		{
			name:  "Nested expressions",
			input: "(1 + 2) * 3",
			want: expression.NewBinary(
				expression.NewGrouping(
					expression.NewBinary(
						expression.NewLiteral(float64(1)),
						token.Token{Type: token.PLUS, Lexeme: "+"},
						expression.NewLiteral(float64(2)),
					),
				),
				token.Token{Type: token.STAR, Lexeme: "*"},
				expression.NewLiteral(float64(3)),
			),
			wantErr: false,
		},
		{
			name:  "Unary expression",
			input: "-5",
			want: expression.NewUnary(
				token.Token{Type: token.MINUS, Lexeme: "-"},
				expression.NewLiteral(float64(5)),
			),
			wantErr: false,
		},
		{
			name:    "Unmatched parenthesis",
			input:   "(1 + 2",
			wantErr: true,
			errMsg:  "Expect ')' after expression",
		},
		{
			name:  "Ternary expression",
			input: "true ? 1 : 2",
			want: expression.NewTernary(
				expression.NewLiteral(true),
				expression.NewLiteral(float64(1)),
				expression.NewLiteral(float64(2)),
			),
			wantErr: false,
		},
		{
			name:    "Invalid ternary",
			input:   "true ? 1",
			wantErr: true,
			errMsg:  "Expect ':' in ternary expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.NewScanner(tt.input)
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatalf("scanner error: %v", err)
			}

			p := NewParser(tokens)
			got, err := p.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Parser.Parse() error = %v, want error message %v", err, tt.errMsg)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Primary(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    expression.Expr
		wantErr bool
	}{
		{
			name:    "Number",
			input:   "42",
			want:    expression.NewLiteral(float64(42)),
			wantErr: false,
		},
		{
			name:    "String",
			input:   "\"hello\"",
			want:    expression.NewLiteral("hello"),
			wantErr: false,
		},
		{
			name:    "True",
			input:   "true",
			want:    expression.NewLiteral(true),
			wantErr: false,
		},
		{
			name:    "False",
			input:   "false",
			want:    expression.NewLiteral(false),
			wantErr: false,
		},
		{
			name:    "Nil",
			input:   "nil",
			want:    expression.NewLiteral(nil),
			wantErr: false,
		},
		{
			name:    "Grouped expression",
			input:   "(42)",
			want:    expression.NewGrouping(expression.NewLiteral(float64(42))),
			wantErr: false,
		},
		{
			name:    "Invalid token",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.NewScanner(tt.input)
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatalf("scanner error: %v", err)
			}

			p := NewParser(tokens)
			got, err := p.primary()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.primary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.primary() = %v, want %v", got, tt.want)
			}
		})
	}
}
