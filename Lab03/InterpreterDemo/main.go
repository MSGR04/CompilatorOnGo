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
var total = 1;
var step = 2;
var label = "factorial";

while (step <= 5) {
    total = total * step;
    print total;
    step = step + 1;
}

var ok = total == 120;
print ok;

if (ok && true) {
    print label + " ok";
}
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

	if warns := analyzer.Warnings(); len(warns) > 0 {
		fmt.Println("Semantic analysis warnings:")
		for _, warning := range warns {
			fmt.Printf("- %s\n", warning)
		}
	}

	fmt.Println("Program output:")

	interpreter := Interpreter.New(os.Stdout)
	if err := interpreter.Interpret(statements); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
	}
}
