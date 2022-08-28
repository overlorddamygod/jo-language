package parser

import (
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/jo/lexer"
)

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
	name           string
	Type           string
	Value          string
	NumericalValue float64
}

func NewLiteralValue(Type, value string) *LiteralValue {
	litVal := LiteralValue{
		name:  "LiteralValue",
		Type:  Type,
		Value: value,
	}
	if Type == L.INT || Type == L.FLOAT {
		litVal.NumericalValue, _ = strconv.ParseFloat(litVal.Value, 32)
	}
	return &litVal
}

func (l *LiteralValue) IsNumber() bool {
	return l.Type == L.INT || l.Type == L.FLOAT
}

func (l *LiteralValue) IsString() bool {
	return l.Type == L.STRING
}

func (l *LiteralValue) IsBoolean() bool {
	return l.Type == L.BOOLEAN
}

func (l *LiteralValue) GetNumber() float64 {
	if l.IsBoolean() {
		if l.Value == "true" {
			return 1
		}
		return 0
	}

	if l.IsString() {
		return 1
	}
	return l.NumericalValue
}

func (l *LiteralValue) GetBoolean() bool {
	if l.IsNumber() {
		return l.GetNumber() > 0
	}
	if l.IsString() {
		return true
	}
	return l.Value == "true"
}

func (l *LiteralValue) GetString() string {
	return l.Value
}

func (l *LiteralValue) NodeName() string {
	return l.name
}
func (l *LiteralValue) Print() {
	fmt.Println(*l)
}

func BooleanLiteral(boolean bool) LiteralValue {
	return *NewLiteralValue(L.BOOLEAN, fmt.Sprintf("%v", boolean))
}
func NumberLiteral(val float64) LiteralValue {
	return *NewLiteralValue(L.FLOAT, fmt.Sprintf("%f", val))
}

func StringLiteral(val string) LiteralValue {
	return *NewLiteralValue(L.STRING, val)
}

type Identifier struct {
	name  string
	Type  string
	Value string
	Token *L.Token
}

func NewIdentifier(value string, token *L.Token) *Identifier {
	return &Identifier{
		name:  "Identifier",
		Type:  "Identifier",
		Value: value,
		Token: token,
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
