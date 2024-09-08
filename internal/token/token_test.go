package token

import "testing"

func TestToken(t *testing.T) {

	tests := []struct {
		tokenType TokenType
		lexeme    string
		literal   interface{}
		line      int
	}{
		{MINUS, "-", nil, 1},
		{NUMBER, "123", 123.0, 2},
		{STRING, "hello", "hello", 3},
	}

	for _, tt := range tests {
		token := NewToken(tt.tokenType, tt.lexeme, tt.literal, tt.line)
		if token.Type != tt.tokenType {
			t.Errorf("Expected token type %v, got %v", tt.tokenType, token.Type)
		}
		if token.Lexeme != tt.lexeme {
			t.Errorf("Expected lexeme %v, got %v", tt.lexeme, token.Lexeme)
		}
		if token.Literal != tt.literal {
			t.Errorf("Expected literal %v, got %v", tt.literal, token.Literal)
		}
		if token.Line != tt.line {
			t.Errorf("Expected line %v, got %v", tt.line, token.Line)
		}

	}
}
