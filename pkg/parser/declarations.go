package parser

import (
	"fmt"
)

type VarDeclStatement struct {
	name       string
	Identifier Node
	Expression *Node
}

func NewVarDeclStatement(identifier, expression Node) *VarDeclStatement {
	return &VarDeclStatement{
		name:       "VarDecl",
		Identifier: identifier,
		Expression: &expression,
	}
}

func (a *VarDeclStatement) NodeName() string {
	return a.name
}

func (a *VarDeclStatement) Print() {
	fmt.Println(a.name)
	a.Identifier.Print()
	(*a.Expression).Print()

	// fmt.Println("Params")
	// for _, p := range a.Params {
	// 	p.Print()
	// }
	// fmt.Println("Body")
	// for _, p := range a.Body.Nodes {
	// 	p.Print()
	// }
}

func (a VarDeclStatement) Type() string {
	return a.name
}

func (v VarDeclStatement) GetLine() int {
	return v.Identifier.GetLine()
}

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

func (v FunctionDeclStatement) GetLine() int {
	return v.Identifier.GetLine()
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

func (b Block) GetLine() int {
	if len(b.Nodes) == 0 {
		return 1
	}
	return b.Nodes[0].GetLine()
}

type StructDeclStatement struct {
	name       string
	Identifier Node
	// Attributes []VarDeclStatement
	Methods []FunctionDeclStatement
}

func NewStructDeclStatement(identifier Node, methods []FunctionDeclStatement) *StructDeclStatement {
	return &StructDeclStatement{
		name:       "StructDecl",
		Identifier: identifier,
		// Attributes: attr,
		Methods: methods,
	}
}

func (a *StructDeclStatement) NodeName() string {
	return a.name
}

func (a *StructDeclStatement) Print() {
	fmt.Println(a.name)
	a.Identifier.Print()

	fmt.Println("Methods")
	for _, p := range a.Methods {
		p.Print()
	}
	// fmt.Println("Body")
	// for _, p := range a.Body.Nodes {
	// 	p.Print()
	// }
}

func (s StructDeclStatement) GetLine() int {
	return s.Identifier.GetLine()
}
