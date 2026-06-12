package main

import (
	"fmt"
	"os"

	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser"
	"CompilatorOnGo/Core/Semantic"
)

func main() {
	source := `
func add(a, b) {
    return a + b;
}

func twice(value) {
    return add(value, value);
}

var result = twice(21);
print result;
`

	fmt.Println("Source program:")
	fmt.Println(source)

	lexer := Lexer.NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
		return
	}

	parser := Parser.New(tokens)
	statements, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
		return
	}

	analyzer := Semantic.NewSemanticAnalyzer()
	analyzer.Analyze(statements)

	if errs := analyzer.Errors(); len(errs) > 0 {
		fmt.Println("Semantic analysis errors:")
		for _, semanticErr := range errs {
			fmt.Printf("- %s\n", semanticErr)
		}
		return
	}

	if warnings := analyzer.Warnings(); len(warnings) > 0 {
		fmt.Println("Semantic analysis warnings:")
		for _, warning := range warnings {
			fmt.Printf("- %s\n", warning)
		}
	}

	fmt.Println("Function declarations and calls were parsed and checked successfully.")
}
