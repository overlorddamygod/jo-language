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

type ReturnStatement struct {
	name       string
	Expression Node
}

func NewReturnStatement(expression Node) *ReturnStatement {
	return &ReturnStatement{
		name:       "ReturnStatement",
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

type BreakStatement struct {
	name string
}

func NewBreakStatement() *BreakStatement {
	return &BreakStatement{
		name: "BreakStatement",
	}
}

func (a *BreakStatement) NodeName() string {
	return a.name
}

func (a *BreakStatement) Print() {
	fmt.Println(a.name)
	// a.Expression.Print()
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

type ForStatement struct {
	name       string
	Initial    Node
	Condition  Node
	Expression Node
	Block      []Node
}

func NewForStatement(Initial, Condition Node, Expression Node, Block []Node) *ForStatement {
	return &ForStatement{
		name:       "FOR",
		Initial:    Initial,
		Condition:  Condition,
		Expression: Expression,
		Block:      Block,
		// ElseBlock: Else,
	}
}

func (a *ForStatement) NodeName() string {
	return a.name
}

func (a *ForStatement) Print() {
	fmt.Println(a.name)
	// a.Condition.Print()

	// fmt.Println("IF")
	// for _, i := range a.IfBlock {
	// 	i.Print()
	// }
	// fmt.Println("ELSE")

	// for _, i := range a.ElseBlock {
	// 	i.Print()
	// }
	// fmt.Println("END IF")
	// a.Identifier.Print()
	// a.Expression.Print()
}
