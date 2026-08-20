// Command astdump prints the AST of one Go expression.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: astdump 'total < 5000'")
		os.Exit(2)
	}

	expression, err := parser.ParseExpr(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse expression: %v\n", err)
		os.Exit(1)
	}

	if err := ast.Print(token.NewFileSet(), expression); err != nil {
		fmt.Fprintf(os.Stderr, "print AST: %v\n", err)
		os.Exit(1)
	}
}
