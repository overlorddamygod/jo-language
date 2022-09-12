package node

import (
	"fmt"

	L "github.com/overlorddamygod/jo/pkg/lexer"
)

type ArrayDecl struct {
	name    string
	Bracket *L.Token
	Values  []Node
}

func NewArrayDecl(node *L.Token, values []Node) *ArrayDecl {
	return &ArrayDecl{
		name:    ARRAY,
		Bracket: node,
		Values:  values,
	}
}

func (a *ArrayDecl) NodeName() string {
	return a.name
}

func (a *ArrayDecl) Print() {
	fmt.Println(a.name)

	fmt.Println("Values")
	for _, p := range a.Values {
		p.Print()
	}

}

func (a ArrayDecl) Type() string {
	return a.name
}

func (v ArrayDecl) GetLine() int {
	return v.Bracket.GetLine()
}
