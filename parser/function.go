package parser

import L "github.com/overlorddamygod/jo/lexer"

func (p *Parser) ParamOrArguments(leftRightParser func() (Node, error), midConditionFunc func(*L.Token) bool) ([]Node, error) {
	var arguments []Node = make([]Node, 0)

	argument, err := leftRightParser()
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, argument)
	for {
		op_token, err := p.lexer.PeekToken(0)
		if err != nil {
			break
		}
		if midConditionFunc(op_token) {
			p.lexer.NextToken()

			argument, err := leftRightParser()

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, argument)
		} else {
			break
		}
	}
	return arguments, nil
}

func (p *Parser) arguments() ([]Node, error) {
	return p.ParamOrArguments(p.expression, func(t *L.Token) bool {
		return t.Type == L.PUNCTUATION && t.Literal == L.COMMA
	})
}
func (p *Parser) parameters() ([]Node, error) {
	return p.ParamOrArguments(p.identifier, func(t *L.Token) bool {
		return t.Type == L.PUNCTUATION && t.Literal == L.COMMA
	})
}
