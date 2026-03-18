package main

import (
	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Parser"
	"CompilatorOnGo/Core/Semantic"
	"fmt"
	"os"
	// compilerlabs "CompilatorOnGo/Core"
)

func main() {
	randomCode := "var x = 10; var y = x + 5; print x;"
	fmt.Println(randomCode)

	lexer := Lexer.NewLexer(randomCode)
	tokens, err := lexer.Tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
		return
	}

	parser := Parser.New(tokens)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
		return
	}

	fmt.Printf("Parsed successfully: %d statements at the top level.\n", len(ast))

	printer := AstPrinter{}
	printer.Print(ast)

	semanticAnalyzer := Semantic.NewSemanticAnalyzer()
	semanticAnalyzer.Analyze(ast)

	errs := semanticAnalyzer.Errors()
	warns := semanticAnalyzer.Warnings()

	if len(errs) > 0 {
		fmt.Println("Semantic analysis errors:")
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
	}

	if len(warns) > 0 {
		fmt.Println("Semantic analysis warnings:")
		for _, warn := range warns {
			fmt.Printf("- %s\n", warn)
		}
	}

	if len(errs) == 0 {
		fmt.Println("Semantic analysis completed successfully (no errors).")
	}
}
