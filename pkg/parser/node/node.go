package node

type Node interface {
	NodeName() string
	Print()
	GetLine() int
}
