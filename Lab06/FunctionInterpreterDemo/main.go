package main

import (
	"fmt"
	"os"

	"CompilatorOnGo/Core/Interpreter"
	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser"
	"CompilatorOnGo/Core/Semantic"
)

func main() {
	source := `
func factorial(n) {
    if (n <= 1) {
        return 1;
    }

    return n * factorial(n - 1);
}

func sumTo(limit) {
    var current = 1;
    var total = 0;

    while (current <= limit) {
        total = total + current;
        current = current + 1;
    }

    return total;
}

print factorial(5);
print sumTo(10);
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

	fmt.Println("Program output:")

	interpreter := Interpreter.New(os.Stdout)
	if err := interpreter.Interpret(statements); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
	}
}
