package parser

import (
	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) functionDecl() (*node.FunctionDeclStatement, error) {
	_, err := p.match(L.KEYWORD, "fn")

	if err != nil {
		return nil, err
	}

	identifier, err := p.identifier()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.LPAREN)

	if err != nil {
		return nil, err
	}

	rightPar, _ := p.lexer.PeekToken(0)

	if rightPar.Type == L.PUNCTUATION && rightPar.Literal == L.RPAREN {
		p.lexer.NextToken()

		block, err := p.block()

		if err != nil {
			return nil, err
		}

		return node.NewFunctionDeclStatement(identifier, []node.Node{}, block), nil
	}

	parameters, err := p.parameters()

	if err != nil {
		return nil, err
	}

	// fmt.Println(arguments)

	_, err = p.match(L.PUNCTUATION, L.RPAREN)

	if err != nil {
		return nil, err
	}

	block, err := p.block()

	if err != nil {
		return nil, err
	}
	return node.NewFunctionDeclStatement(identifier, parameters, block), nil
}

func (p *Parser) call() (node.Node, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		leftParenOrDot, _ := p.lexer.PeekToken(0)
		if leftParenOrDot.Type == L.PUNCTUATION && leftParenOrDot.Literal == L.LPAREN {
			p.lexer.NextToken()

			rightPar, _ := p.lexer.PeekToken(0)

			if rightPar.Type == L.PUNCTUATION && rightPar.Literal == L.RPAREN {
				p.lexer.NextToken()

				expr = node.NewFunctionCall(expr, []node.Node{})
				continue
			}

			arguments, err := p.arguments()

			if err != nil {
				return nil, err
			}

			// fmt.Println(arguments)

			rightParen, err := p.lexer.PeekToken(0)

			if err != nil || !(rightParen.Type == L.PUNCTUATION && rightParen.Literal == L.RPAREN) {
				return nil, JoError.New(p.lexer, rightParen, JoError.SyntaxError, "Expected )")
			}
			p.lexer.NextToken()
			expr = node.NewFunctionCall(expr, arguments)
		} else if leftParenOrDot.Type == L.OPERATOR && leftParenOrDot.Literal == L.FULL_STOP {
			p.lexer.NextToken()

			iden, err := p.identifier()

			if err != nil {
				return nil, err
			}

			expr = node.NewGetExpr(iden, expr)
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) _return() (node.Node, error) {
	returnToken, err := p.match(L.KEYWORD, "return")

	if err != nil {
		return nil, err
	}

	semicolon, _ := p.lexer.PeekToken(0)

	if semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON {
		return node.NewReturnStatement(returnToken, nil), nil
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	return node.NewReturnStatement(returnToken, exp), nil
}

func (p *Parser) ParamOrArguments(leftRightParser func() (node.Node, error), midConditionFunc func(*L.Token) bool) ([]node.Node, error) {
	var arguments []node.Node = make([]node.Node, 0)

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

func (p *Parser) arguments() ([]node.Node, error) {
	return p.ParamOrArguments(p.expression, func(t *L.Token) bool {
		return t.Type == L.PUNCTUATION && t.Literal == L.COMMA
	})
}
func (p *Parser) parameters() ([]node.Node, error) {
	return p.ParamOrArguments(p.identifier, func(t *L.Token) bool {
		return t.Type == L.PUNCTUATION && t.Literal == L.COMMA
	})
}
