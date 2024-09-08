package environment

import (
	"fmt"
	Token "interpreter/internal/token"
)

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}), enclosing: enclosing}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(token Token.Token) interface{} {
	if val, exists := e.values[token.Lexeme]; exists {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.Get(token)
	}
	panic(fmt.Sprintf("Undefined variable %s ", token.Lexeme))
}

func (e *Environment) Assign(token Token.Token, value interface{}) {
	if _, ok := e.values[token.Lexeme]; ok {
		e.values[token.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(token, value)
	}

	panic(fmt.Sprintf("Undefined variable %s ", token.Lexeme))

}
