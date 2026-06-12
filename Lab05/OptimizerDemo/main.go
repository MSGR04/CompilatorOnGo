package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"CompilatorOnGo/Core/Interpreter"
	"CompilatorOnGo/Core/Lexer"
	"CompilatorOnGo/Core/Optimizer"
	"CompilatorOnGo/Core/Parser"
	"CompilatorOnGo/Core/Parser/Ast"
	"CompilatorOnGo/Core/Semantic"
)


func printAST(statements []Ast.Statement) string {
	var sb strings.Builder
	sb.WriteString("Program\n")
	for i, stmt := range statements {
		printStmt(&sb, stmt, "", i == len(statements)-1)
	}
	return sb.String()
}

func printStmt(sb *strings.Builder, node interface{}, indent string, isLast bool) {
	if node == nil {
		return
	}
	marker := "├── "
	childIndent := indent + "│   "
	if isLast {
		marker = "└── "
		childIndent = indent + "    "
	}

	switch n := node.(type) {
	case *Ast.VarStatement:
		fmt.Fprintf(sb, "%s%sVarStatement: %s\n", indent, marker, n.Name)
		if n.Initializer != nil {
			printExpr(sb, n.Initializer, childIndent, true)
		}
	case *Ast.PrintStatement:
		fmt.Fprintf(sb, "%s%sPrintStatement\n", indent, marker)
		printExpr(sb, n.Expr, childIndent, true)
	case *Ast.ExpressionStatement:
		fmt.Fprintf(sb, "%s%sExpressionStatement\n", indent, marker)
		printExpr(sb, n.Expr, childIndent, true)
	case *Ast.IfStatement:
		fmt.Fprintf(sb, "%s%sIfStatement\n", indent, marker)
		printExpr(sb, n.Condition, childIndent, false)
		printStmt(sb, n.ThenBranch, childIndent, n.ElseBranch == nil)
		if n.ElseBranch != nil {
			printStmt(sb, n.ElseBranch, childIndent, true)
		}
	case *Ast.WhileStatement:
		fmt.Fprintf(sb, "%s%sWhileStatement\n", indent, marker)
		printExpr(sb, n.Condition, childIndent, false)
		printStmt(sb, n.Body, childIndent, true)
	case *Ast.BlockStatement:
		fmt.Fprintf(sb, "%s%sBlockStatement\n", indent, marker)
		for i, s := range n.Statements {
			printStmt(sb, s, childIndent, i == len(n.Statements)-1)
		}
	case *Ast.FunctionStatement:
		fmt.Fprintf(sb, "%s%sFunctionStatement: %s(%s)\n", indent, marker, n.Name, strings.Join(n.Params, ", "))
		for i, s := range n.Body {
			printStmt(sb, s, childIndent, i == len(n.Body)-1)
		}
	case *Ast.ReturnStatement:
		fmt.Fprintf(sb, "%s%sReturnStatement\n", indent, marker)
		if n.Value != nil {
			printExpr(sb, n.Value, childIndent, true)
		}
	default:
		fmt.Fprintf(sb, "%s%s<unknown stmt %T>\n", indent, marker, node)
	}
}

func printExpr(sb *strings.Builder, node interface{}, indent string, isLast bool) {
	marker := "├── "
	childIndent := indent + "│   "
	if isLast {
		marker = "└── "
		childIndent = indent + "    "
	}

	switch n := node.(type) {
	case *Ast.NumberExpression:
		fmt.Fprintf(sb, "%s%sNumber(%g)\n", indent, marker, n.Value)
	case *Ast.StringExpression:
		fmt.Fprintf(sb, "%s%sString(%q)\n", indent, marker, n.Value)
	case *Ast.BooleanExpression:
		fmt.Fprintf(sb, "%s%sBoolean(%v)\n", indent, marker, n.Value)
	case *Ast.VariableExpression:
		fmt.Fprintf(sb, "%s%sVariable(%s)\n", indent, marker, n.Name)
	case *Ast.BinaryExpression:
		fmt.Fprintf(sb, "%s%sBinary(%s)\n", indent, marker, tokenName(n.Operator))
		printExpr(sb, n.Left, childIndent, false)
		printExpr(sb, n.Right, childIndent, true)
	case *Ast.UnaryExpression:
		fmt.Fprintf(sb, "%s%sUnary(%s)\n", indent, marker, tokenName(n.Operator))
		printExpr(sb, n.Right, childIndent, true)
	case *Ast.AssignExpression:
		fmt.Fprintf(sb, "%s%sAssign(%s)\n", indent, marker, n.Name)
		printExpr(sb, n.Value, childIndent, true)
	case *Ast.CallExpression:
		fmt.Fprintf(sb, "%s%sCall\n", indent, marker)
		printExpr(sb, n.Callee, childIndent, len(n.Arguments) == 0)
		for i, arg := range n.Arguments {
			printExpr(sb, arg, childIndent, i == len(n.Arguments)-1)
		}
	default:
		fmt.Fprintf(sb, "%s%s<unknown expr %T>\n", indent, marker, node)
	}
}

func tokenName(t Lexer.TokenType) string {
	names := map[Lexer.TokenType]string{
		Lexer.PLUS: "+", Lexer.MINUS: "-", Lexer.STAR: "*", Lexer.SLASH: "/",
		Lexer.EQ: "=", Lexer.EQEQ: "==", Lexer.NEQ: "!=",
		Lexer.LT: "<", Lexer.GT: ">", Lexer.LTEQ: "<=", Lexer.GTEQ: ">=",
		Lexer.AND: "&&", Lexer.OR: "||", Lexer.EXCL: "!",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}


func parse(source string) ([]Ast.Statement, error) {
	lexer := Lexer.NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("lexer: %w", err)
	}
	stmts, err := Parser.New(tokens).Parse()
	if err != nil {
		return nil, fmt.Errorf("parser: %w", err)
	}
	return stmts, nil
}

func checkSemantics(stmts []Ast.Statement) error {
	analyzer := Semantic.NewSemanticAnalyzer()
	analyzer.Analyze(stmts)
	if errs := analyzer.Errors(); len(errs) > 0 {
		return fmt.Errorf("semantic errors: %v", errs)
	}
	return nil
}

func runAndCapture(stmts []Ast.Statement) (string, error) {
	var buf bytes.Buffer
	interp := Interpreter.New(&buf)
	if err := interp.Interpret(stmts); err != nil {
		return "", err
	}
	return buf.String(), nil
}


type testCase struct {
	name   string
	source string
}

func runTest(tc testCase) bool {
	sep := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n", sep)
	fmt.Printf("TEST: %s\n", tc.name)
	fmt.Printf("Source:\n%s\n", tc.source)

	before, err := parse(tc.source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  FAIL (parse): %v\n", err)
		return false
	}

	if err := checkSemantics(before); err != nil {
		fmt.Fprintf(os.Stderr, "  FAIL (semantic): %v\n", err)
		return false
	}

	opt := Optimizer.New()
	after := opt.Optimize(before)

	beforeStr := printAST(before)
	afterStr := printAST(after)

	fmt.Println("── AST before optimization:")
	fmt.Print(beforeStr)
	fmt.Println("── AST after optimization:")
	fmt.Print(afterStr)

	if beforeStr == afterStr {
		fmt.Println("  (no change — no constant sub-expressions found)")
	} else {
		fmt.Println("  Optimization applied.")
	}

	if err := checkSemantics(after); err != nil {
		fmt.Printf("  FAIL: optimized AST failed semantic check: %v\n", err)
		return false
	}

	outBefore, errB := runAndCapture(before)
	outAfter, errA := runAndCapture(after)

	bothErr := errB != nil && errA != nil
	if errB != nil && errA == nil {
		fmt.Printf("  FAIL: original errored (%v) but optimized did not\n", errB)
		return false
	}
	if errB == nil && errA != nil {
		fmt.Printf("  FAIL: original succeeded but optimized errored (%v)\n", errA)
		return false
	}
	if bothErr {
		fmt.Printf("  Both produced runtime error (expected): %v\n", errB)
		fmt.Println("  PASS: semantics preserved (same error)")
		return true
	}
	if outBefore != outAfter {
		fmt.Printf("  FAIL: output mismatch\n    before: %q\n    after:  %q\n", outBefore, outAfter)
		return false
	}
	if outBefore != "" {
		fmt.Printf("  Output: %s", outBefore)
	}
	fmt.Println("  PASS: outputs match")
	return true
}


func main() {
	tests := []testCase{
		{
			name:   "Number arithmetic: 2 + 3 * 4",
			source: `var x = 2 + 3 * 4; print x;`,
		},
		{
			name:   "Chained arithmetic: (10 - 4) / 2 + 1",
			source: `var x = (10 - 4) / 2 + 1; print x;`,
		},
		{
			name:   "Unary minus on literal",
			source: `var x = -7 + 2; print x;`,
		},
		{
			name:   "String concatenation of literals",
			source: `var s = "hello" + ", " + "world"; print s;`,
		},
		{
			name:   "Boolean comparison: 5 > 3",
			source: `var b = 5 > 3; print b;`,
		},
		{
			name:   "Logical AND/OR on literals",
			source: `var a = true && false; var b = false || true; print a; print b;`,
		},
		{
			name:   "Unary NOT on literal",
			source: `var b = !false; print b;`,
		},
		{
			name:   "Mixed: constant cond in if",
			source: `if (2 < 5) { print "yes"; } else { print "no"; }`,
		},
		{
			name:   "No folding: variable operand",
			source: `var x = 10; var y = x + 5; print y;`,
		},
		{
			name:   "Semantics preserved: division by zero NOT folded",
			source: `var x = 5 / 0; print x;`,
		},
		{
			name:   "Constant in function body",
			source: `func double(n) { return n + 0; } print double(21);`,
		},
		{
			name:   "String concat mixed with variable (no fold)",
			source: `var name = "Alice"; var greeting = "Hello, " + name; print greeting;`,
		},
	}

	passed := 0
	for _, tc := range tests {
		if runTest(tc) {
			passed++
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("═", 60))
	fmt.Printf("Results: %d / %d passed\n", passed, len(tests))
}
