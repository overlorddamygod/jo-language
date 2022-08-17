package parser

import "fmt"

// type Node

type BinaryExpression struct {
	name  string
	op    string
	left  Node
	right Node
}

func NewBinaryExpression(op string, left, right Node) *BinaryExpression {
	return &BinaryExpression{
		name:  "BinaryExpression",
		op:    op,
		left:  left,
		right: right,
	}
}

func (b *BinaryExpression) NodeName() string {
	return b.name
}

func (b *BinaryExpression) Print() {
	fmt.Println(b.name)
	fmt.Println(b.op)
	b.left.Print()
	b.right.Print()
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

type Node interface {
	NodeName() string
	Print()
}
