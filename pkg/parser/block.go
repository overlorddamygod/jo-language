package parser

import (
	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) block() (*node.Block, error) {
	leftCurly, err := p.lexer.NextToken()

	// fmt.Println("LEFTCURLY", leftCurly)
	if err != nil || !(leftCurly.Type == L.PUNCTUATION && leftCurly.Literal == L.LBRACE) {
		return nil, JoError.New(p.lexer, leftCurly, JoError.SyntaxError, "Expected {")
	}

	// fmt.Println("BLOCK")

	block, err := p.declarations()

	if err != nil {
		return nil, err
	}

	rightCurly, err := p.lexer.NextToken()

	if err != nil || !(rightCurly.Type == L.PUNCTUATION && rightCurly.Literal == L.RBRACE) {
		return nil, JoError.New(p.lexer, rightCurly, JoError.SyntaxError, "Expected }")
	}

	return node.NewBlock(block), nil
}
