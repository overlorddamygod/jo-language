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

type IfStatement struct {
	name      string
	Condition Node
	IfBlock   []Node
	ElseBlock []Node
}

func NewIfStatement(Condition Node, If []Node) *IfStatement {
	return &IfStatement{
		name:      "IF",
		Condition: Condition,
		IfBlock:   If,
		// ElseBlock: Else,
	}
}

func (i *IfStatement) Else(Else []Node) *IfStatement {
	i.ElseBlock = Else
	return i
}

func (a *IfStatement) NodeName() string {
	return a.name
}

func (a *IfStatement) Print() {
	fmt.Println(a.name)
	a.Condition.Print()

	fmt.Println("IF")
	for _, i := range a.IfBlock {
		i.Print()
	}
	fmt.Println("ELSE")

	for _, i := range a.ElseBlock {
		i.Print()
	}
	fmt.Println("END IF")
	// a.Identifier.Print()
	// a.Expression.Print()
}
