package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]

	defineAst(outputDir, "Expr", []string{
		"Assign   : Name Token.Token, Value Expr",
		"Binary   : Left Expr, Operator Token.Token, Right Expr",
		"Ternary   : Condition Expr, TrueExpression Expr, FalseExpression Expr",
		"Grouping : Expr Expr",
		"Literal  : Value interface{}",
		"Unary    : Operator Token.Token, Right Expr",
		"Variable : Name Token.Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Expression:  expr Expr",
		"Print: Expression Expr",
		"Var:  Name Token.Token, Initializer Expr",
		"Block: Statements Stmt[]",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	path := outputDir + "/" + baseName + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintln(file, "package expression")
	//Change this path to your token package definition
	fmt.Fprintln(file, `import Token "interpreter/internal/token"`)
	fmt.Fprintln(file)
	defineVisitor(file, baseName, types)
	fmt.Fprintln(file, "type", baseName, "interface{")
	fmt.Fprintf(file, "	Accept(visitor %sVisitor) interface{}\n", baseName)
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])
		defineType(file, baseName, className, fields)

	}
}

func defineType(file *os.File, baseName, className, fieldList string) {
	fmt.Fprintf(file, "type %s struct {\n", className)
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		fmt.Fprintf(file, "    %s\n", field)
	}
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)
	// constructor
	fmt.Fprintf(file, "func New%s(%s) *%s {\n", className, fieldList, className)
	fmt.Fprintf(file, "    return &%s{\n", className)
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		fmt.Fprintf(file, "        %s: %s,\n", name, name)
	}
	fmt.Fprintln(file, "    }")
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

	// visitor
	fmt.Fprintf(file, "func (e *%s) Accept(visitor %sVisitor) interface{} {\n", className, baseName)
	fmt.Fprintf(file, "    return visitor.Visit%s%s(e)\n", className, baseName)
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)
}

func defineVisitor(file *os.File, baseName string, types []string) {
	fmt.Fprintf(file, "type %sVisitor interface {\n", baseName)
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintf(file, "    Visit%s%s(%s *%s) interface{}\n", typeName, baseName, strings.ToLower(baseName), typeName)
	}
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)

}
