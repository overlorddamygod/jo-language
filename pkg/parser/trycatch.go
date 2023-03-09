package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) tryCatch() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "try")

	if err != nil {
		return nil, err
	}

	tryBlock, err := p.block()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.KEYWORD, "catch")

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.LPAREN)

	if err != nil {
		return nil, err
	}

	catchVar, err := p.identifier()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.RPAREN)

	if err != nil {
		return nil, err
	}

	catchBlock, err := p.block()

	if err != nil {
		return nil, err
	}

	catchId := catchVar.(*node.Identifier)

	return node.NewTryCatchStatement(tryBlock, catchId, catchBlock), nil
}

func (p *Parser) throw() (node.Node, error) {
	throw, err := p.match(L.KEYWORD, "throw")

	if err != nil {
		return nil, err
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	return node.NewThrowStatement(throw, exp), nil
}
