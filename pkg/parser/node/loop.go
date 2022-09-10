package node

import (
	"fmt"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

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
	}
}

func (a *ForStatement) NodeName() string {
	return a.name
}

func (a *ForStatement) Print() {
	fmt.Println(a.name)
	a.Initial.Print()
	a.Condition.Print()
	fmt.Println(a.Block)
	a.Expression.Print()
}
func (a ForStatement) GetLine() int {
	return a.Initial.GetLine()
}

type WhileStatement struct {
	name      string
	Condition Node
	Block     *Block
}

func NewWhileStatement(Condition Node, block *Block) *WhileStatement {
	return &WhileStatement{
		name:      "WHILE",
		Condition: Condition,
		Block:     block,
	}
}

func (a *WhileStatement) NodeName() string {
	return a.name
}

func (a *WhileStatement) Print() {
	fmt.Println(a.name)
	a.Condition.Print()
	fmt.Println(a.Block)
}
func (a WhileStatement) GetLine() int {
	return a.Condition.GetLine()
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
