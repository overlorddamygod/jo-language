package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) switchDecl() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "switch")

	if err != nil {
		return nil, err
	}

	exp, err := p.condition()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.LBRACE)

	if err != nil {
		return nil, err
	}

	switchStmt := node.NewSwitchStatement(exp)

	var cases []node.Case = make([]node.Case, 0)

	for {
		err = p.peekMatch(0, L.KEYWORD, "case")

		if err != nil {
			break
		}

		case_, err := p._Case()

		if err != nil {
			return nil, err
		}

		cases = append(cases, *case_)
	}

	switchStmt.SetCases(cases)

	for {
		err = p.peekMatch(0, L.KEYWORD, "default")
		// fmt.Println("SAD", p.lexer.PeekToken(0))
		if err != nil {
			break
		}

		p.lexer.NextToken()

		_, err = p.match(L.PUNCTUATION, L.COLON)

		if err != nil {
			return nil, err
		}

		block, err := p.switchBlock()
		if err != nil {
			return nil, err
		}

		switchStmt.SetDefault(*node.NewBlock(block))
		break
	}

	_, err = p.match(L.PUNCTUATION, L.RBRACE)

	if err != nil {
		return nil, err
	}

	return switchStmt, nil
}

func (p *Parser) caseValues() ([]node.Node, error) {
	return p.ParamOrArguments(p.expression, func(t *L.Token) bool {
		return t.Type == L.PUNCTUATION && t.Literal == L.COMMA
	})
}

func (p *Parser) _Case() (*node.Case, error) {
	_, err := p.match(L.KEYWORD, "case")

	if err != nil {
		return nil, err
	}

	values, err := p.caseValues()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.COLON)

	if err != nil {
		return nil, err
	}
	block, err := p.switchBlock()

	if err != nil {
		return nil, err
	}

	return node.NewCase(values, *node.NewBlock(block)), nil
}

func (p *Parser) switchBlock() ([]node.Node, error) {
	var declarations []node.Node = make([]node.Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if (token.Type == L.PUNCTUATION && token.Literal == L.RBRACE) || (token.Type == L.KEYWORD && token.Literal == "case") || (token.Type == L.KEYWORD && token.Literal == "default") {
			break
		}
		dec, err := p.declaration()
		if err != nil {
			return declarations, err
		}
		declarations = append(declarations, dec)
	}
	return declarations, nil
}
