package parser

import "fmt"

// type Node

type BinaryExpression struct {
	name  string
	Op    string
	Left  Node
	Right Node
}

func NewBinaryExpression(op string, left, right Node) *BinaryExpression {
	return &BinaryExpression{
		name:  "BinaryExpression",
		Op:    op,
		Left:  left,
		Right: right,
	}
}

func (b *BinaryExpression) NodeName() string {
	return b.name
}

func (b *BinaryExpression) Print() {
	fmt.Println(b.name)
	fmt.Println(b.Op)
	b.Left.Print()
	b.Right.Print()
}

type LiteralValue struct {
	name  string
	Type  string
	Value string
}

func NewLiteralValue(Type, value string) *LiteralValue {
	return &LiteralValue{
		name:  "LiteralValue",
		Type:  Type,
		Value: value,
	}
}

func (l *LiteralValue) NodeName() string {
	return l.name
}
func (l *LiteralValue) Print() {
	fmt.Println(*l)
}

type Identifier struct {
	name  string
	Type  string
	Value string
}

func NewIdentifier(value string) *Identifier {
	return &Identifier{
		name:  "Identifier",
		Type:  "Identifier",
		Value: value,
	}
}

func (i *Identifier) NodeName() string {
	return i.name
}
func (i *Identifier) Print() {
	fmt.Println(*i)
}

type Node interface {
	NodeName() string
	Print()
}
