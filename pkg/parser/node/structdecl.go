package node

import "fmt"

type StructDeclStatement struct {
	name       string
	Identifier Node
	// Attributes []VarDeclStatement
	Methods []FunctionDeclStatement
}

func NewStructDeclStatement(identifier Node, methods []FunctionDeclStatement) *StructDeclStatement {
	return &StructDeclStatement{
		name:       STRUCT_DECL,
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
