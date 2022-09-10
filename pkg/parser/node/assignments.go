package node

import "fmt"

type AssignmentStatement struct {
	name       string
	Identifier Node
	Op         string
	Expression Node
}

func NewAssignmentStatement(identifier Node, op string, exp Node) *AssignmentStatement {
	return &AssignmentStatement{
		name:       ASSIGNMENT,
		Identifier: identifier,
		Op:         op,
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

func (a AssignmentStatement) GetLine() int {
	return a.Identifier.GetLine()
}
