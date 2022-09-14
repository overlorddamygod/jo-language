package node

import "fmt"

type SwitchStatement struct {
	name    string
	Test    Node
	Cases   []Case
	Default Block
}

type Case struct {
	Values []Node
	Block  Block
}

func NewCase(Values []Node, Block Block) *Case {
	return &Case{
		Values: Values,
		Block:  Block,
	}
}

func NewSwitchStatement(Test Node) *SwitchStatement {
	return &SwitchStatement{
		name: SWITCH,
		Test: Test,
	}
}

func (a *SwitchStatement) SetCases(cases []Case) {
	a.Cases = cases
}
func (a *SwitchStatement) SetDefault(defaultBlock Block) {
	a.Default = defaultBlock
}

func (a *SwitchStatement) NodeName() string {
	return a.name
}

func (a SwitchStatement) GetLine() int {
	return a.Test.GetLine()
}

func (a *SwitchStatement) Print() {
	fmt.Println(a.name)
}
