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
		fmt.Println(err)
		return
		// panic(err)
	}
	// fmt.Println(tokens)

	parser := Parser.NewParser(lexer)

	node := parser.Parse()

	// for _, s := range node {
	// 	s.Print()
	// }

	evaluator := Eval.NewEvaluator(node)

	evaluator.Eval()
}
