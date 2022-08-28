package parser

import (
	"errors"

	L "github.com/overlorddamygod/jo/lexer"
)

type Parser struct {
	lexer *L.Lexer
}

func NewParser(lexer *L.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() ([]Node, error) {
	var statements []Node = make([]Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == L.EOF {
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

func (p *Parser) Statements() ([]Node, error) {
	var statements []Node = make([]Node, 0)
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

func (p *Parser) statement() (Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "if":
		return p.ifElse()
	case "for":
		return p.For()
	}
	second, _ := p.lexer.PeekToken(1)

	if second.Literal == L.ASSIGN {
		assignment, err := p.assignment()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(assignment)
	}

	functioncall, err := p.functionCall()

	if err != nil {
		return nil, errors.New("expected ;")
	}
	return p.matchSemicolon(functioncall)
}

func (p *Parser) matchSemicolon(node Node) (Node, error) {
	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, errors.New("expected ;")
	}
	return node, nil
}

func (p *Parser) functionCall() (Node, error) {
	functionName, err := p.identifier()

	if err != nil {
		return nil, errors.New("expected identifier")
	}

	leftParenthesis, err := p.lexer.NextToken()

	if err != nil || !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		return nil, errors.New("expected (")
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	rightParenthesis, err := p.lexer.NextToken()

	if err != nil || !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		return nil, errors.New("expected )")
	}

	return NewFunctionCallStatement(functionName, exp), nil
}

func (p *Parser) ifElse() (Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil || identifier.Literal != "if" {
		return nil, errors.New("expected if")

	}

	leftParenthesis, err := p.lexer.NextToken()

	if err != nil || !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		return nil, errors.New("expected (")
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	rightParenthesis, err := p.lexer.NextToken()

	if err != nil || !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		return nil, errors.New("expected )")

	}
	// fmt.Println("HERE")

	ifBlock, err := p.block()

	if err != nil {
		return nil, err
	}

	ifStatement := NewIfStatement(exp, ifBlock)

	token, _ := p.lexer.PeekToken(0)

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

func (p *Parser) For() (Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil || identifier.Literal != "for" {
		return nil, errors.New("EXPECTED for")
	}

	leftParenthesis, err := p.lexer.NextToken()

	if err != nil || !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		return nil, errors.New("EXPECTED (")
	}

	assignment, err := p.assignment()

	if err != nil {
		return nil, err
	}

	// fmt.Println(assignment)
	// assignment.Print()

	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, errors.New("EXPECTED ;")
	}

	condition, err := p.expression()

	if err != nil {
		return nil, errors.New("expected looping condition")
	}

	// condition.Print()

	semicolon, err = p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, errors.New("EXPECTED ;")
	}

	exp, err := p.assignment()

	if err != nil {
		return nil, err
	}
	// exp.Print()

	rightParenthesis, err := p.lexer.NextToken()

	if err != nil || !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		return nil, errors.New("EXPECTED )")
	}
	// fmt.Println("HERE")

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return NewForStatement(assignment, condition, exp, block), nil
}

func (p *Parser) block() ([]Node, error) {
	leftCurly, err := p.lexer.NextToken()

	// fmt.Println("LEFTCURLY", leftCurly)
	if err != nil || !(leftCurly.Type == L.PUNCTUATION && leftCurly.Literal == L.LBRACE) {
		return nil, errors.New("expected {")
	}

	// fmt.Println("BLOCK")

	block, err := p.Statements()

	if err != nil {
		return nil, err
	}

	rightCurly, err := p.lexer.NextToken()

	if err != nil || !(rightCurly.Type == L.PUNCTUATION && rightCurly.Literal == L.RBRACE) {
		return nil, errors.New("expected }")
	}

	return block, nil
}

func (p *Parser) identifier() (Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil {
		return nil, err
	}
	return NewIdentifier(identifier.Literal), nil
}

func (p *Parser) assignment() (Node, error) {
	identifier, err := p.identifier()

	if err != nil {
		return nil, errors.New("EXPECTED identifier")
	}

	equals, err := p.lexer.NextToken()

	if err != nil || !(equals.Type == L.OPERATOR && equals.Literal == "=") {
		return nil, errors.New("EXPECTED =")
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	// semicolon, err := p.lexer.NextToken()

	// if err != nil {

	// }

	// if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
	// 	fmt.Println("EXPECTED ;")
	// 	return nil
	// }

	return NewAssignmentStatement(identifier, exp), nil
}

// parse expression
func (p *Parser) expression() (Node, error) {
	return p.logicOr()
}

func (p *Parser) binary(leftRightParser func() (Node, error), midConditionFunc func(*L.Token) bool) (Node, error) {
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

			left = NewBinaryExpression(string(op_token.Literal), left, right)
		} else {
			break
		}
	}
	return left, nil
}

func (p *Parser) logicOr() (Node, error) {
	return p.binary(p.logicAnd, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.OR)
	})
}

func (p *Parser) logicAnd() (Node, error) {
	return p.binary(p.equality, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.AND)
	})
}
func (p *Parser) equality() (Node, error) {
	return p.binary(p.comparison, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.NOT_EQ || t.Literal == L.EQ)
	})
}
func (p *Parser) comparison() (Node, error) {
	return p.binary(p.term, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.LT || t.Literal == L.LT_EQ || t.Literal == L.GT || t.Literal == L.GT_EQ)
	})
}

func (p *Parser) term() (Node, error) {
	return p.binary(p.factor, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.PLUS || t.Literal == L.MINUS)
	})
}

func (p *Parser) factor() (Node, error) {
	return p.binary(p.primary, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.SLASH || t.Literal == L.ASTERISK || t.Literal == L.PERCENT)
	})
}

func (p *Parser) primary() (Node, error) {
	token, err := p.lexer.NextToken()

	if err != nil {
		return nil, err
	}

	if token.Type == L.STRING || token.Type == L.INT || token.Type == L.FLOAT {
		// token.Print()
		return NewLiteralValue(string(token.Type), token.Literal), nil
	}

	if token.Type == L.KEYWORD && (token.Literal == L.TRUE || token.Literal == L.FALSE) {
		// token.Print()
		return NewLiteralValue("BOOLEAN", token.Literal), nil
	}

	if token.Type == L.IDENTIFIER {
		// token.Print()
		return NewIdentifier(token.Literal), nil
	}

	if token.Type == L.PUNCTUATION && token.Literal == L.LPAREN {
		// token.Print()

		e, err := p.expression()

		if err != nil {
			return nil, err
		}
		token, err := p.lexer.NextToken()

		if err != nil || !(token.Type == L.PUNCTUATION && token.Literal == L.RPAREN) {
			return nil, errors.New("expected )")

		}
		return e, nil
	}
	return nil, errors.New("unexpected primary value")
}
