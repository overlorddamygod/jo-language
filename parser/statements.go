package parser

import "fmt"

type AssignmentStatement struct {
	name       string
	Identifier Node
	Expression Node
}

func NewAssignmentStatement(identifier Node, exp Node) *AssignmentStatement {
	return &AssignmentStatement{
		name:       "ASSIGNMENT",
		Identifier: identifier,
		Expression: exp,
	}
}

func (a *AssignmentStatement) NodeName() string {
	return a.name
}

func (a *AssignmentStatement) Print() {
	fmt.Println(a.name)
	a.Identifier.Print()
	a.Expression.Print()
}

type FunctionCallStatement struct {
	name       string
	Identifier Node
	Expression Node
}

func NewFunctionCallStatement(identifier Node, exp Node) *FunctionCallStatement {
	return &FunctionCallStatement{
		name:       "FUNCTION_CALL",
		Identifier: identifier,
		Expression: exp,
	}
}

func (a *FunctionCallStatement) NodeName() string {
	return a.name
}

func (a *FunctionCallStatement) Print() {
	fmt.Println(a.name)
	a.Identifier.Print()
	a.Expression.Print()
}
