package main

import (
	"bufio"
	"fmt"
	"os"

	"CompilatorOnGo/Core"
	"CompilatorOnGo/Core/Lexer"
)

func main() {
	generator := compilerlabs.NewRandomProgramGenerator(0)

	randomCode := generator.Generate(10)

	fmt.Println("=== СГЕНЕРИРОВАННЫЙ КОД ===")
	fmt.Println(randomCode)
	fmt.Println("===========================")
	fmt.Println()

	lexer := Lexer.NewLexer(randomCode)
	tokens, err := lexer.Tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== ТОКЕНЫ ===")
	for _, token := range tokens {
		fmt.Println(token)
	}

	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
