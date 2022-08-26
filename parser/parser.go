package parser

import (
	"errors"
	"fmt"

	"github.com/overlorddamygod/lexer/lexer"
	L "github.com/overlorddamygod/lexer/lexer"
)

type Parser struct {
	lexer *L.Lexer
}

func NewParser(lexer *L.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() []Node {
	var statements []Node = make([]Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == lexer.EOF {
			break
		}
		st, _ := p.statement()
		statements = append(statements, st)
	}
	return statements
}

func (p *Parser) Statements() []Node {
	var statements []Node = make([]Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == lexer.PUNCTUATION && token.Literal == L.RBRACE {
			break
		}
		st, _ := p.statement()
		statements = append(statements, st)
	}
	return statements
}

func (p *Parser) statement() (Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "if":
		return p.ifElse(), nil
	case "for":
		return p.For(), nil
	}
	second, _ := p.lexer.PeekToken(1)

	if second.Literal == lexer.ASSIGN {
		return p.matchSemicolon(p.assignment())
	}

	return p.matchSemicolon(p.functionCall())
}

func (p *Parser) matchSemicolon(node Node) (Node, error) {
	semicolon, err := p.lexer.NextToken()

	if err != nil {

	}

	if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		fmt.Println("EXPECTED ;")
		return nil, errors.New("expected ;")
	}
	return node, nil
}

func (p *Parser) functionCall() Node {
	functionName, _ := p.identifier()

	leftParenthesis, _ := p.lexer.NextToken()

	if !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		fmt.Println("EXPECTED (")
		return nil
	}

	exp := p.expression()

	if exp == nil {
		fmt.Println("EXPECTED expression")
		return nil
	}

	rightParenthesis, _ := p.lexer.NextToken()

	if !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		fmt.Println("EXPECTED )")
		return nil
	}

	return NewFunctionCallStatement(functionName, exp)
}

func (p *Parser) ifElse() Node {
	identifier, _ := p.lexer.NextToken()

	if identifier.Literal != "if" {
		fmt.Println("EXPECTED if")
		return nil
	}

	leftParenthesis, _ := p.lexer.NextToken()

	if !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		fmt.Println("EXPECTED (")
		return nil
	}

	exp := p.expression()

	rightParenthesis, _ := p.lexer.NextToken()

	if !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		fmt.Println("EXPECTED )")
		return nil
	}
	// fmt.Println("HERE")

	ifBlock := p.block()

	ifStatement := NewIfStatement(exp, ifBlock)

	token, _ := p.lexer.PeekToken(0)

	if token.Literal == "else" {
		p.lexer.NextToken()
		elseBlock := p.block()
		ifStatement.Else(elseBlock)
	}
	return ifStatement
}

func (p *Parser) For() Node {
	identifier, _ := p.lexer.NextToken()

	if identifier.Literal != "for" {
		fmt.Println("EXPECTED for")
		return nil
	}

	leftParenthesis, _ := p.lexer.NextToken()

	if !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		fmt.Println("EXPECTED (")
		return nil
	}

	assignment := p.assignment()

	// fmt.Println(assignment)
	// assignment.Print()

	semicolon, err := p.lexer.NextToken()

	// fmt.Println(semicolon)
	if err != nil {

	}

	if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		fmt.Println("EXPECTED ;")
		return nil
	}

	condition := p.expression()

	// condition.Print()

	semicolon, err = p.lexer.NextToken()
	// fmt.Println("HERE", semicolon)
	if err != nil {

	}

	if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		fmt.Println("EXPECTED ;")
		return nil
	}

	exp := p.assignment()

	// exp.Print()

	rightParenthesis, _ := p.lexer.NextToken()

	if !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		fmt.Println("EXPECTED )")
		return nil
	}
	// fmt.Println("HERE")

	block := p.block()

	return NewForStatement(assignment, condition, exp, block)
}

func (p *Parser) block() []Node {
	leftCurly, _ := p.lexer.NextToken()

	// fmt.Println("LEFTCURLY", leftCurly)
	if !(leftCurly.Type == L.PUNCTUATION && leftCurly.Literal == L.LBRACE) {
		fmt.Println("EXPECTED {")
		return nil
	}

	// fmt.Println("BLOCK")

	block := p.Statements()

	rightCurly, _ := p.lexer.NextToken()

	if !(rightCurly.Type == L.PUNCTUATION && rightCurly.Literal == L.RBRACE) {
		fmt.Println("EXPECTED }")
		return nil
	}

	return block
}

func (p *Parser) identifier() (Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil {
		return nil, err
	}
	return NewIdentifier(identifier.Literal), nil
}

func (p *Parser) assignment() Node {
	identifier, err := p.identifier()

	if err != nil {
		fmt.Println("EXPECTED identifier")
		return nil
	}

	equals, err := p.lexer.NextToken()

	if err != nil {

	}

	if !(equals.Type == L.OPERATOR && equals.Literal == "=") {
		fmt.Println("EXPECTED =")
		return nil
	}

	exp := p.expression()

	if exp == nil {
		fmt.Println("EXPECTED expression")
		return nil
	}

	// semicolon, err := p.lexer.NextToken()

	// if err != nil {

	// }

	// if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
	// 	fmt.Println("EXPECTED ;")
	// 	return nil
	// }

	return NewAssignmentStatement(identifier, exp)
}

// parse expression
func (p *Parser) expression() Node {
	return p.logicOr()
}

func (p *Parser) binary(leftRightParser func() Node, midConditionFunc func(*L.Token) bool) Node {
	left := leftRightParser()
	// if err != nil {

	// }
	for {
		op_token, err := p.lexer.PeekToken(0)
		if err != nil {
			break
		}
		if midConditionFunc(op_token) {
			p.lexer.NextToken()

			right := leftRightParser()

			left = NewBinaryExpression(string(op_token.Literal), left, right)
		} else {
			break
		}
	}
	return left
}

func (p *Parser) logicOr() Node {
	return p.binary(p.logicAnd, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.OR)
	})
}

func (p *Parser) logicAnd() Node {
	return p.binary(p.equality, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.AND)
	})
}
func (p *Parser) equality() Node {
	return p.binary(p.comparison, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.NOT_EQ || t.Literal == L.EQ)
	})
}
func (p *Parser) comparison() Node {
	return p.binary(p.term, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.LT || t.Literal == L.LT_EQ || t.Literal == L.GT || t.Literal == L.GT_EQ)
	})
}

func (p *Parser) term() Node {
	return p.binary(p.factor, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.PLUS || t.Literal == L.MINUS)
	})
}

func (p *Parser) factor() Node {
	return p.binary(p.primary, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.SLASH || t.Literal == L.ASTERISK || t.Literal == L.PERCENT)
	})
}

func (p *Parser) primary() Node {
	token, err := p.lexer.NextToken()

	if err != nil {

	}

	if token.Type == L.STRING || token.Type == L.INT || token.Type == L.FLOAT {
		// token.Print()
		return NewLiteralValue(string(token.Type), token.Literal)
	}

	if token.Type == L.KEYWORD && (token.Literal == L.TRUE || token.Literal == L.FALSE) {
		// token.Print()
		return NewLiteralValue("BOOLEAN", token.Literal)
	}

	if token.Type == L.IDENTIFIER {
		// token.Print()
		return NewIdentifier(token.Literal)
	}

	if token.Type == L.PUNCTUATION && token.Literal == L.LPAREN {
		// token.Print()

		e := p.expression()
		token, err := p.lexer.NextToken()

		if err != nil {

		}
		if token.Type == L.PUNCTUATION && token.Literal == L.RPAREN {
			token.Print()
			return e
		} else {
			fmt.Println("Expected )")
			return nil
		}
	}
	return nil
}
