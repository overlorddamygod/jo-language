package parser

import (
	"strings"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

// parse expression
func (p *Parser) expression() (node.Node, error) {
	return p.assignment()
}

func (p *Parser) assignment() (node.Node, error) {
	exp, err := p.logicOr()

	// fmt.Println("ORRR", exp)
	pos := p.lexer.GetTokenPos()
	op, err := p.matchMany(L.OPERATOR, L.ASSIGN, L.PLUS_ASSIGN, L.MINUS_ASSIGN, L.ASTERISK_ASSIGN, L.SLASH_ASSIGN, L.BANG_ASSIGN, L.PIPE_ASSIGN, L.AND_ASSIGN, L.OR_ASSIGN, L.AMPERSAND_ASSIGN, L.PERCENT_ASSIGN)
	if err != nil {
		p.lexer.SetTokenPos(pos)
		return exp, nil
	}
	opLiteral := getOpFromAssignment(op.Literal)
	// fmt.Println()
	ass, err := p.assignment()

	if err != nil {
		return ass, err
	}

	return node.NewAssignmentStatement(exp, opLiteral, ass), nil
}

func (p *Parser) binary(leftRightParser func() (node.Node, error), midConditionFunc func(*L.Token) bool) (node.Node, error) {
	left, err := leftRightParser()
	if err != nil {
		return nil, err
	}
	for {
		op_token, err := p.lexer.PeekToken(0)
		if err != nil {
			break
		}
		if midConditionFunc(op_token) {
			p.lexer.NextToken()

			right, err := leftRightParser()

			if err != nil {
				return nil, err
			}

			left = node.NewBinaryExpression(string(op_token.Literal), left, right)
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) logicOr() (node.Node, error) {
	return p.binary(p.logicAnd, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.OR)
	})
}

func (p *Parser) logicAnd() (node.Node, error) {
	return p.binary(p.equality, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.AND)
	})
}
func (p *Parser) equality() (node.Node, error) {
	return p.binary(p.comparison, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.NOT_EQ || t.Literal == L.EQ)
	})
}
func (p *Parser) comparison() (node.Node, error) {
	return p.binary(p.term, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.LT || t.Literal == L.LT_EQ || t.Literal == L.GT || t.Literal == L.GT_EQ)
	})
}

func (p *Parser) term() (node.Node, error) {
	return p.binary(p.factor, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.PLUS || t.Literal == L.MINUS)
	})
}

func (p *Parser) factor() (node.Node, error) {
	return p.binary(p.unary, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.SLASH || t.Literal == L.ASTERISK || t.Literal == L.PERCENT)
	})
}

func (p *Parser) unary() (node.Node, error) {
	token, _ := p.lexer.PeekToken(0)

	if token.Type == L.OPERATOR && (token.Literal == L.BANG || token.Literal == L.UNARY_PLUS || token.Literal == L.UNARY_MINUS) {
		p.lexer.NextToken()
		unary_, _ := p.unary()
		return node.NewUnaryExpression(string(token.Literal), unary_, token), nil
	}

	return p.call()
}

func (p *Parser) primary() (node.Node, error) {
	token, err := p.lexer.NextToken()

	// fmt.Println("BOOO", token)

	if err != nil {
		return nil, JoError.New(p.lexer, token, JoError.SyntaxError, "Expected primary value")
	}

	if token.Type == L.STRING || token.Type == L.INT || token.Type == L.FLOAT {
		// token.Print()
		return node.NewLiteralValue(string(token.Type), token.Literal), nil
	}

	if token.Type == L.KEYWORD && (token.Literal == L.TRUE || token.Literal == L.FALSE) {
		// token.Print()
		return node.NewLiteralValue("BOOLEAN", token.Literal), nil
	}

	if token.Type == L.IDENTIFIER {
		// token.Print()
		return node.NewIdentifier(token.Literal, token), nil
	}
	// fmt.Println("HERE")

	if token.Type == L.PUNCTUATION && token.Literal == L.LPAREN {
		// token.Print()

		e, err := p.expression()

		if err != nil {
			return nil, err
		}
		token, err := p.lexer.NextToken()

		if err != nil || !(token.Type == L.PUNCTUATION && token.Literal == L.RPAREN) {
			return nil, JoError.New(p.lexer, token, JoError.SyntaxError, "Expected )")

		}
		return e, nil
	}
	return nil, JoError.New(p.lexer, token, JoError.SyntaxError, "unexpected value")
}

func (p *Parser) identifier() (node.Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil || identifier.Type == L.EOF || identifier.Type != L.IDENTIFIER {
		return nil, JoError.New(p.lexer, identifier, JoError.SyntaxError, "Expected identifier")
	}

	if L.IsKeyword(identifier.Literal) {
		return nil, JoError.New(p.lexer, identifier, JoError.SyntaxError, "Variable name cannot be a keyword")
	}

	return node.NewIdentifier(identifier.Literal, identifier), nil
}

func getOpFromAssignment(op string) string {
	if op == "=" {
		return op
	}
	return strings.Split(op, "=")[0]
}
