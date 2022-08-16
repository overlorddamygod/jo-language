package main

import (
	"fmt"
	"os"

	Lexer "github.com/overlorddamygod/lexer/lexer"
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

		panic(err)
	}
	fmt.Println(tokens)
}
