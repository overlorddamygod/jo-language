package node

import "fmt"

type IfStatement struct {
	name      string
	IfBlocks  []*ConditionBlock
	ElseBlock *Block
}

func NewIfStatement(If []*ConditionBlock) *IfStatement {
	return &IfStatement{
		name:     IF,
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
		name:      CONDITION_BLOCK,
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
