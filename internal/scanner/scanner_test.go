package scanner

import (
	"interpreter/internal/token"
	"reflect"
	"testing"
)

func TestScanner_ScanTokens(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []token.Token
		wantErr bool
	}{
		{
			name:  "Empty input",
			input: "",
			want:  []token.Token{{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1}},
		},
		{
			name:  "Single character tokens",
			input: "(){},.-+;*",
			want: []token.Token{
				{Type: token.LEFT_PAREN, Lexeme: "(", Line: 1},
				{Type: token.RIGHT_PAREN, Lexeme: ")", Line: 1},
				{Type: token.LEFT_BRACE, Lexeme: "{", Line: 1},
				{Type: token.RIGHT_BRACE, Lexeme: "}", Line: 1},
				{Type: token.COMMA, Lexeme: ",", Line: 1},
				{Type: token.DOT, Lexeme: ".", Line: 1},
				{Type: token.MINUS, Lexeme: "-", Line: 1},
				{Type: token.PLUS, Lexeme: "+", Line: 1},
				{Type: token.SEMICOLON, Lexeme: ";", Line: 1},
				{Type: token.STAR, Lexeme: "*", Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "One or two character tokens",
			input: "! != = == < <= > >=",
			want: []token.Token{
				{Type: token.BANG, Lexeme: "!", Line: 1},
				{Type: token.BANG_EQUAL, Lexeme: "!=", Line: 1},
				{Type: token.EQUAL, Lexeme: "=", Line: 1},
				{Type: token.EQUAL_EQUAL, Lexeme: "==", Line: 1},
				{Type: token.LESS, Lexeme: "<", Line: 1},
				{Type: token.LESS_EQUAL, Lexeme: "<=", Line: 1},
				{Type: token.GREATER, Lexeme: ">", Line: 1},
				{Type: token.GREATER_EQUAL, Lexeme: ">=", Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Comments",
			input: "// This is a comment\n5",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "5", Literal: float64(5), Line: 2},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 2},
			},
		},
		{
			name:  "Strings",
			input: "\"Hello, World!\"",
			want: []token.Token{
				{Type: token.STRING, Lexeme: "\"Hello, World!\"", Literal: "Hello, World!", Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Numbers",
			input: "123 45.67",
			want: []token.Token{
				{Type: token.NUMBER, Lexeme: "123", Literal: float64(123), Line: 1},
				{Type: token.NUMBER, Lexeme: "45.67", Literal: 45.67, Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Keywords and identifiers",
			input: "var language = \"next\";",
			want: []token.Token{
				{Type: token.VAR, Lexeme: "var", Line: 1},
				{Type: token.IDENTIFIER, Lexeme: "language", Line: 1},
				{Type: token.EQUAL, Lexeme: "=", Line: 1},
				{Type: token.STRING, Lexeme: "\"next\"", Literal: "next", Line: 1},
				{Type: token.SEMICOLON, Lexeme: ";", Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:  "Ternary operator",
			input: "true ? 1 : 2",
			want: []token.Token{
				{Type: token.TRUE, Lexeme: "true", Line: 1},
				{Type: token.QUESTION_MARK, Lexeme: "?", Line: 1},
				{Type: token.NUMBER, Lexeme: "1", Literal: float64(1), Line: 1},
				{Type: token.COLON, Lexeme: ":", Line: 1},
				{Type: token.NUMBER, Lexeme: "2", Literal: float64(2), Line: 1},
				{Type: token.EOF, Lexeme: "", Literal: nil, Line: 1},
			},
		},
		{
			name:    "Unterminated string",
			input:   "\"Unterminated",
			wantErr: true,
		},
		{
			name:    "Invalid character",
			input:   "@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.input)
			got, err := s.ScanTokens()

			if (err != nil) != tt.wantErr {
				t.Errorf("Scanner.ScanTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.ScanTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_isAtEnd(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{"Empty string", "", true},
		{"Non-empty string", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.source)
			if got := s.isAtEnd(); got != tt.want {
				t.Errorf("Scanner.isAtEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_advance(t *testing.T) {
	s := NewScanner("ab")
	if got := s.advance(); got != 'a' {
		t.Errorf("Scanner.advance() = %v, want %v", got, 'a')
	}
	if got := s.advance(); got != 'b' {
		t.Errorf("Scanner.advance() = %v, want %v", got, 'b')
	}
}

func TestScanner_peek(t *testing.T) {
	s := NewScanner("ab")
	if got := s.peek(); got != 'a' {
		t.Errorf("Scanner.peek() = %v, want %v", got, 'a')
	}
	s.advance()
	if got := s.peek(); got != 'b' {
		t.Errorf("Scanner.peek() = %v, want %v", got, 'b')
	}
}

func TestScanner_peekNext(t *testing.T) {
	s := NewScanner("abc")
	if got := s.peekNext(); got != 'b' {
		t.Errorf("Scanner.peekNext() = %v, want %v", got, 'b')
	}
	s.advance()
	if got := s.peekNext(); got != 'c' {
		t.Errorf("Scanner.peekNext() = %v, want %v", got, 'c')
	}
}
