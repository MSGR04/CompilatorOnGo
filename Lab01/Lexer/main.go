package main

import (
	"CompilatorOnGo/Core/Lexer"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	codeExample := "var x = 123; print x + 5;"

	lexer := Lexer.NewLexer(codeExample)
	tokens, err := lexer.Tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Tokens for example program ---")
	for _, token := range tokens {
		fmt.Println(token)
	}

	for i := 0; i < 3; i++ {
		fmt.Println()
		fmt.Println("--- Generating random test program ---")
		randomProgram := generateRandomTestProgram()
		fmt.Println(randomProgram)

		randomLexer := Lexer.NewLexer(randomProgram)
		randomTokens, err := randomLexer.Tokenize()
		if err != nil {
			fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
			os.Exit(1)
		}

		for _, token := range randomTokens {
			fmt.Println(token)
		}
	}

	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func generateRandomTestProgram() string {
	rand.Seed(time.Now().UnixNano())
	variables := []string{"a", "b", "c", "x", "y", "z"}
	operators := []string{"+", "-", "*", "/"}

	var program string

	for i := 0; i < 5; i++ {
		varName := variables[rand.Intn(len(variables))]
		number := rand.Intn(99) + 1 // 1..99
		program += fmt.Sprintf("var %s = %d;\n", varName, number)
	}

	for i := 0; i < 5; i++ {
		var1 := variables[rand.Intn(len(variables))]
		var2 := variables[rand.Intn(len(variables))]
		op := operators[rand.Intn(len(operators))]
		program += fmt.Sprintf("print %s %s %s;\n", var1, op, var2)
	}

	return program
}
