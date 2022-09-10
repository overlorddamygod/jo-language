package node

import "fmt"

type VarDeclStatement struct {
	name       string
	Identifier Node
	Expression *Node
}

func NewVarDeclStatement(identifier, expression Node) *VarDeclStatement {
	return &VarDeclStatement{
		name:       VAR_DECL,
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
}

func (a VarDeclStatement) Type() string {
	return a.name
}

func (v VarDeclStatement) GetLine() int {
	return v.Identifier.GetLine()
}
