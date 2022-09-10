package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) block() (*node.Block, error) {
	_, err := p.match(L.PUNCTUATION, L.LBRACE)

	if err != nil {
		return nil, err
	}

	block, err := p.declarations()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.RBRACE)

	if err != nil {
		return nil, err
	}

	return node.NewBlock(block), nil
}
