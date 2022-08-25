package main

import (
	"fmt"
	"os"

	Eval "github.com/overlorddamygod/lexer/eval"
	Lexer "github.com/overlorddamygod/lexer/lexer"
	Parser "github.com/overlorddamygod/lexer/parser"
)

func main() {
	dat, err := os.ReadFile("./example.jo")

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

	evaluator := Eval.NewEvaluator(node)

	evaluator.Eval()
}
