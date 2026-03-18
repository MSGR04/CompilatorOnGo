package main

import (
	"CompilatorOnGo/Core/Lexer"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	codeExample := "var x; print x + 5;"

	lexer := Lexer.NewLexer(codeExample)
	tokens, err := lexer.Tokenize()
	if err != nil {
		fmt.Println("Lexer error:", err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	fmt.Print("\nPress Enter to exit...")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}

func GenerateRandomTestProgram() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	variables := []string{"a", "b", "c", "x", "y", "z"}
	operators := []string{"+", "-", "*", "/"}

	var b strings.Builder

	for i := 0; i < 5; i++ {
		varName := variables[rng.Intn(len(variables))]
		number := rng.Intn(99) + 1 // 1..99 как Next(1,100)
		fmt.Fprintf(&b, "var %s = %d;\n", varName, number)
	}

	for i := 0; i < 5; i++ {
		var1 := variables[rng.Intn(len(variables))]
		var2 := variables[rng.Intn(len(variables))]
		op := operators[rng.Intn(len(operators))]
		fmt.Fprintf(&b, "print %s %s %s;\n", var1, op, var2)
	}

	return b.String()
}
