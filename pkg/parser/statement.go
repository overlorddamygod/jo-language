package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) statements() ([]node.Node, error) {
	var statements []node.Node = make([]node.Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == L.PUNCTUATION && token.Literal == L.RBRACE {
			break
		}
		st, err := p.statement()
		if err != nil {
			return statements, err
		}
		statements = append(statements, st)
	}
	return statements, nil
}
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
	// TODO: ADD BLOCK STATEMENT
	// case "{":
	// return p.block()
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
