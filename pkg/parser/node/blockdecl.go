package node

import "fmt"

type Block struct {
	name  string
	Nodes []Node
}

func NewBlock(nodes []Node) *Block {
	return &Block{
		name:  BLOCK,
		Nodes: nodes,
	}
}

func (b Block) Print() {
	fmt.Println(b.name)
	for _, node := range b.Nodes {
		node.Print()
	}
}

func (b Block) NodeName() string {
	return b.name
}

func (b Block) GetLine() int {
	if len(b.Nodes) == 0 {
		return 1
	}
	return b.Nodes[0].GetLine()
}
