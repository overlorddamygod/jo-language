package node

type Block struct {
	name  string
	Nodes []Node
}

func NewBlock(nodes []Node) *Block {
	return &Block{
		name:  "Block",
		Nodes: nodes,
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
