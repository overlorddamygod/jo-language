package parser

import (
	"fmt"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

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

func (a AssignmentStatement) GetLine() int {
	return a.Identifier.GetLine()
}

type ReturnStatement struct {
	name       string
	token      lexer.Token
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

func (a ReturnStatement) GetLine() int {
	return a.token.GetLine()
}

type BreakStatement struct {
	name  string
	token lexer.Token
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

func (a BreakStatement) GetLine() int {
	return a.token.GetLine()
}

type ContinueStatement struct {
	name  string
	token lexer.Token
}

func NewContinueStatement() *ContinueStatement {
	return &ContinueStatement{
		name: "ContinueStatement",
	}
}

func (c *ContinueStatement) NodeName() string {
	return c.name
}

func (c *ContinueStatement) Print() {
	fmt.Println(c.name)
	// a.Expression.Print()
}

func (a ContinueStatement) GetLine() int {
	return a.token.GetLine()
}

type IfStatement struct {
	name      string
	IfBlocks  []*ConditionBlock
	ElseBlock *Block
}

func NewIfStatement(If []*ConditionBlock) *IfStatement {
	return &IfStatement{
		name:     "IF",
		IfBlocks: If,
		// ElseBlock: Else,
	}
}

func (i *IfStatement) HasIfs() bool {
	return i.IfBlocks != nil && len(i.IfBlocks) > 0
}

func (i *IfStatement) Else(Else *Block) *IfStatement {
	i.ElseBlock = Else
	return i
}

func (i *IfStatement) HasElse() bool {
	return i.ElseBlock != nil && len(i.ElseBlock.Nodes) != 0
}

func (a *IfStatement) NodeName() string {
	return a.name
}

func (a IfStatement) GetLine() int {
	return a.IfBlocks[0].Block.GetLine()
}

func (a *IfStatement) Print() {
	fmt.Println(a.name)
	// a.Print()

	fmt.Println("IF")
	for _, i := range a.IfBlocks {
		for _, j := range i.Block.Nodes {
			j.Print()
		}
	}

	if a.HasElse() {
		fmt.Println("ELSE")
		for _, i := range a.ElseBlock.Nodes {
			i.Print()
		}
	}
	fmt.Println("END IF")
	// a.Identifier.Print()
	// a.Expression.Print()
}

type ConditionBlock struct {
	name      string
	Condition Node
	Block     *Block
}

func NewConditionBlock(Condition Node, Block *Block) *ConditionBlock {
	return &ConditionBlock{
		name:      "ConditionBlock",
		Condition: Condition,
		Block:     Block,
	}
}

func (c ConditionBlock) NodeName() string {
	return c.name
}

func (c ConditionBlock) Print() {
	fmt.Println("CONDITION")
	c.Condition.Print()

	fmt.Println("BLOCK")

	for _, i := range c.Block.Nodes {
		i.Print()
	}
}

func (a ConditionBlock) GetLine() int {
	return a.Block.GetLine()
}

type ForStatement struct {
	name       string
	Initial    Node
	Condition  Node
	Expression Node
	Block      *Block
}

func NewForStatement(Initial, Condition Node, Expression Node, block *Block) *ForStatement {
	return &ForStatement{
		name:       "FOR",
		Initial:    Initial,
		Condition:  Condition,
		Expression: Expression,
		Block:      block,
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
func (a ForStatement) GetLine() int {
	return a.Initial.GetLine()
}
