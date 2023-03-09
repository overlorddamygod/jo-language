package main

import (
	"os"

	Eval "github.com/overlorddamygod/jo/internal/eval"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

func main() {
	if len(os.Args) == 1 {
		stdio.Io.Print("Filename not specified")
		stdio.Io.Print("Usage: jo <filename>")
		return
	}
	file := os.Args[1]

	Eval.Init(file)
}
