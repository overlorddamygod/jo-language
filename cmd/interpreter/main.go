package main

import (
	"os"

	Eval "github.com/overlorddamygod/jo/internal/eval"
	Lexer "github.com/overlorddamygod/jo/pkg/lexer"
	Parser "github.com/overlorddamygod/jo/pkg/parser"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

func main() {
	if len(os.Args) == 1 {
		stdio.Io.Print("Filename not specified")
		stdio.Io.Print("Usage: jo <filename>")
		return
	}
	file := os.Args[1]

	dat, err := os.ReadFile(file)

	if err != nil {
		stdio.Io.Print("Unable to read file", dat)
		return
	}

	lexer := Lexer.NewLexer(string(dat))

	tokens, err := lexer.Lex()
	if err != nil {
		stdio.Io.Print(tokens)
		stdio.Io.Print("[Lexer]\n\n", err)
		return
	}
	// stdio.Io.Print(tokens)

	parser := Parser.NewParser(lexer)

	node, err := parser.Parse()

	// for _, s := range node {
	// 	s.Print()
	// }
	if err != nil {
		stdio.Io.Print("[Parser]\n\n", err)
		return
	}

	// for _, s := range node {
	// 	s.Print()
	// }

	evaluator := Eval.NewEvaluator(lexer, node)

	_, err = evaluator.Eval()

	if err != nil {
		stdio.Io.Printf("[Evaluator]\n\n%s", err)
	}
}
