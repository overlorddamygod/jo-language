package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) condition() (node.Node, error) {
	_, err := p.match(L.PUNCTUATION, L.LPAREN)

	if err != nil {
		return nil, err
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.RPAREN)

	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) ifElse() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "if")

	if err != nil {
		return nil, err
	}

	exp, err := p.condition()

	if err != nil {
		return nil, err
	}

	// fmt.Println("HERE")

	ifBlock, err := p.block()

	if err != nil {
		return nil, err
	}

	var ifs []*node.ConditionBlock = make([]*node.ConditionBlock, 0)
	ifs = append(ifs, node.NewConditionBlock(exp, ifBlock))

	token, _ := p.lexer.PeekToken(0)

	if token.Literal == "elif" {

		for {
			p.lexer.NextToken()

			exp, err := p.expression()

			if err != nil {
				return nil, err
			}

			block, err := p.block()

			if err != nil {
				return nil, err
			}

			ifs = append(ifs, node.NewConditionBlock(exp, block))

			token, _ = p.lexer.PeekToken(0)

			if token.Literal != "elif" {
				break
			}
		}
	}

	ifStatement := node.NewIfStatement(ifs)

	if token.Literal == "else" {
		p.lexer.NextToken()
		elseBlock, err := p.block()

		if err != nil {
			return nil, err
		}
		ifStatement.Else(elseBlock)
	}
	return ifStatement, nil
}

func (p *Parser) _break() (node.Node, error) {
	breakToken, err := p.match(L.KEYWORD, "break")

	if err != nil {
		return nil, err
	}

	return node.NewBreakStatement(breakToken), nil
}

func (p *Parser) _continue() (node.Node, error) {
	continueToken, err := p.match(L.KEYWORD, "continue")

	if err != nil {
		return nil, err
	}

	return node.NewContinueStatement(continueToken), nil
}
