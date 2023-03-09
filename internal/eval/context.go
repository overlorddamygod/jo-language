package eval

import "fmt"

type Context struct {
	name   string
	line   int
	parent *Context
}

func NewContext(name string, line int, parent *Context) *Context {
	return &Context{
		name:   name,
		line:   line,
		parent: parent,
	}
}

func (c *Context) SetLine(line int) {
	c.line = line
}

func (c *Context) Print() {
	fmt.Printf("%d %s\n", c.line, c.name)

	if c.parent != nil {
		c.parent.Print()
	}
}

func (c *Context) Pop() *Context {
	return c.parent
}
