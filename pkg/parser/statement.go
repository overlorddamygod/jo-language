package parser

import (
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) statement() (node.Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "if":
		return p.ifElse()
	case "while":
		return p.While()
	case "for":
		return p.For()
	case "return":
		ret, err := p._return()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(ret)
	case "try":
		return p.tryCatch()
	case "throw":
		throw, err := p.throw()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(throw)
	case "{":
		return p.block()
	case "break":
		br, err := p._break()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(br)

	case "continue":
		con, err := p._continue()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(con)
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	return p.matchSemicolon(exp)
}
