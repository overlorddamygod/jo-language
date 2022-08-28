package main

import (
	"fmt"
	"os"

	Eval "github.com/overlorddamygod/jo/eval"
	Lexer "github.com/overlorddamygod/jo/lexer"
	Parser "github.com/overlorddamygod/jo/parser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Filename not specified")
		fmt.Println("Usage: jo <filename>")
		return
	}
	dat, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	lexer := Lexer.NewLexer(string(dat))

	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Println(tokens)
		fmt.Println("[Lexer]\n", err)
		return
	}
	// fmt.Println(tokens)

	parser := Parser.NewParser(lexer)

	node, err := parser.Parse()

	if err != nil {
		fmt.Println("[Parser]", err)
		return
	}

	// for _, s := range node {
	// 	s.Print()
	// }

	evaluator := Eval.NewEvaluator(node)

	err = evaluator.Eval()

	if err != nil {
		fmt.Printf("[Evaluator] %s", err)
	}
}
