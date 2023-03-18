package node

import (
	"fmt"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

type ImportStatement struct {
	name  string
	File  *lexer.Token
	Alias *lexer.Token
}

func NewImport(file *lexer.Token, alias *lexer.Token) *ImportStatement {
	return &ImportStatement{
		name:  IMPORT,
		File:  file,
		Alias: alias,
	}
}

func (c ImportStatement) NodeName() string {
	return c.name
}

func (c ImportStatement) Print() {
	fmt.Println(c.NodeName(), c.File.Literal)
}

func (a ImportStatement) GetLine() int {
	return a.File.GetLine()
}

type ExportStatement struct {
	name string
	Expr Node
}

func NewExport(expr Node) *ExportStatement {
	return &ExportStatement{
		name: EXPORT,
		Expr: expr,
	}
}

func (c ExportStatement) NodeName() string {
	return c.name
}

func (c ExportStatement) Print() {
	fmt.Println(c.NodeName())
	c.Expr.Print()
}

func (a ExportStatement) GetLine() int {
	return a.Expr.GetLine()
}
