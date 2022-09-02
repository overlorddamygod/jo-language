package parser

import (
	"fmt"
)

type FunctionDeclStatement struct {
	name       string
	Identifier Node
	Params     []Node
	Body       *Block
}

func NewFunctionDeclStatement(identifier Node, params []Node, body *Block) *FunctionDeclStatement {
	return &FunctionDeclStatement{
		name:       "FunctionDecl",
		Identifier: identifier,
		Params:     params,
		Body:       body,
	}
}

func (a *FunctionDeclStatement) NodeName() string {
	return a.name
}

func (a *FunctionDeclStatement) Print() {
	fmt.Println(a.name)
	a.Identifier.Print()

	fmt.Println("Params")
	for _, p := range a.Params {
		p.Print()
	}
	fmt.Println("Body")
	for _, p := range a.Body.Nodes {
		p.Print()
	}
}

func (a FunctionDeclStatement) Type() string {
	return a.name
}

type Block struct {
	name  string
	Nodes []Node
}

func NewBlock(nodes []Node) *Block {
	return &Block{
		name:  "Block",
		Nodes: nodes,
	}
}

func (b Block) NodeName() string {
	return b.name
}
