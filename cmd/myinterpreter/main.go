package main

import (
	"fmt"
	"os"

	"interpreter/internal/interpreter"
	"interpreter/internal/parser"
	scanner "interpreter/internal/scanner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" && command != "evaluate" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	s := scanner.NewScanner(string(fileContents))
	tokens, err := s.ScanTokens()

	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}

	if command == "tokenize" {
		for _, t := range tokens {
			fmt.Printf("%v\n", t)

		}
	}

	p := parser.NewParser(tokens)
	expr, err := p.Parse()

	if err != nil {
		fmt.Println(err)
		os.Exit(70)
	}

	/* if command == "parse" {

		printer := &expression.AstPrinter{}
		fmt.Println(printer.Print(expr))

	} */
	if command == "evaluate" {
		i := interpreter.NewInterpreter()
		i.Interpret(expr)
		/* 	fmt.Println(expr) */
	}

}
