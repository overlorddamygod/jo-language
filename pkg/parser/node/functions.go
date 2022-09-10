package node

import (
	"fmt"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

type FunctionDeclStatement struct {
	name       string
	Identifier Node
	Params     []Node
	Body       *Block
}

func NewFunctionDeclStatement(identifier Node, params []Node, body *Block) *FunctionDeclStatement {
	return &FunctionDeclStatement{
		name:       FUNCTION_DECL,
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

type FunctionCall struct {
	name       string
	Identifier Node
	Arguments  []Node
}

func NewFunctionCall(identifier Node, arguments []Node) *FunctionCall {
	return &FunctionCall{
		name:       FUNCTION_CALL,
		Identifier: identifier,
		Arguments:  arguments,
	}
}

func (b *FunctionCall) NodeName() string {
	return b.name
}

func (b *FunctionCall) Print() {
	fmt.Println(b.name)
	b.Identifier.Print()

	fmt.Println("Arguments")
	for _, s := range b.Arguments {
		s.Print()
	}
}

func (f FunctionCall) GetLine() int {
	return f.Identifier.GetLine()
}

type ReturnStatement struct {
	name       string
	token      lexer.Token
	Expression Node
}

func NewReturnStatement(expression Node) *ReturnStatement {
	return &ReturnStatement{
		name:       RETURN,
		Expression: expression,
	}
}

func (a *ReturnStatement) NodeName() string {
	return a.name
}

func (a *ReturnStatement) Print() {
	fmt.Println(a.name)
	a.Expression.Print()
	// a.Expression.Print()
}

func (a ReturnStatement) GetLine() int {
	return a.token.GetLine()
}
