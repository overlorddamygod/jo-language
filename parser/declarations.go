package parser

import (
	"fmt"
)

type FunctionDeclStatement struct {
	name       string
	Identifier Node
	Params     []Node
	Body       []Node
}

func NewFunctionDeclStatement(identifier Node, params, body []Node) *FunctionDeclStatement {
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
	for _, p := range a.Body {
		p.Print()
	}
}

func (a FunctionDeclStatement) Type() string {
	return a.name
}
