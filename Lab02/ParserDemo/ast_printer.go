package main

import (
	"fmt"

	"CompilatorOnGo/Core/Parser/Ast"
)


type AstPrinter struct{}

func (p *AstPrinter) Print(statements []Ast.Statement) {
	fmt.Println("Root (Program)")
	for i, stmt := range statements {
		isLast := i == len(statements)-1
		p.printNode(stmt, "", isLast)
	}
}

func (p *AstPrinter) printNode(node interface{}, indent string, isLast bool) {
	if node == nil {
		return
	}

	marker := "├── "
	if isLast {
		marker = "└── "
	}
	fmt.Print(indent + marker)

	childIndent := indent
	if isLast {
		childIndent += "    "
	} else {
		childIndent += "│   "
	}

	switch n := node.(type) {
	case *Ast.VarStatement:
		fmt.Printf("VarStatement: %s\n", n.Name)
		if n.Initializer != nil {
			p.printNode(n.Initializer, childIndent, true)
		}

	case *Ast.PrintStatement:
		fmt.Println("PrintStatement")
		p.printNode(n.Expr, childIndent, true)

	case *Ast.IfStatement:
		fmt.Println("IfStatement")
		p.printNode(n.Condition, childIndent, false)
		p.printNode(n.ThenBranch, childIndent, n.ElseBranch == nil)
		if n.ElseBranch != nil {
			p.printNode(n.ElseBranch, childIndent, true)
		}

	case *Ast.WhileStatement:
		fmt.Println("WhileStatement")
		p.printNode(n.Condition, childIndent, false)
		p.printNode(n.Body, childIndent, true)

	case *Ast.BlockStatement:
		fmt.Println("BlockStatement")
		for j, stmt := range n.Statements {
			p.printNode(stmt, childIndent, j == len(n.Statements)-1)
		}

	case *Ast.ExpressionStatement:
		fmt.Println("ExpressionStatement")
		p.printNode(n.Expr, childIndent, true)

	case *Ast.BinaryExpression:
		fmt.Printf("BinaryExpression: %v\n", n.Operator)
		p.printNode(n.Left, childIndent, false)
		p.printNode(n.Right, childIndent, true)

	case *Ast.UnaryExpression:
		fmt.Printf("UnaryExpression: %v\n", n.Operator)
		p.printNode(n.Right, childIndent, true)

	case *Ast.AssignExpression:
		fmt.Printf("AssignExpression: %s =\n", n.Name)
		p.printNode(n.Value, childIndent, true)

	case *Ast.NumberExpression:
		fmt.Printf("Number: %v\n", n.Value)

	case *Ast.VariableExpression:
		fmt.Printf("Variable: %s\n", n.Name)

	default:
		fmt.Printf("Unknown Node: %T\n", node)
	}
}
