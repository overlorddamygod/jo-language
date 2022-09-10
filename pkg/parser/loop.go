package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) For() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "for")

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.LPAREN)

	if err != nil {
		return nil, err
	}

	vardecl, err := p.vardecl()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.SEMICOLON)

	if err != nil {
		return nil, err
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.SEMICOLON)

	if err != nil {
		return nil, err
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}
	// exp.Print()

	_, err = p.match(L.PUNCTUATION, L.RPAREN)

	if err != nil {
		return nil, err
	}

	// fmt.Println("HERE")

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return node.NewForStatement(vardecl, condition, exp, block), nil
}

func (p *Parser) While() (node.Node, error) {
	if _, err := p.match(L.KEYWORD, "while"); err != nil {
		return nil, err
	}

	if _, err := p.match(L.PUNCTUATION, L.LPAREN); err != nil {
		return nil, err
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.match(L.PUNCTUATION, L.RPAREN); err != nil {
		return nil, err
	}

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return node.NewWhileStatement(condition, block), nil
}
