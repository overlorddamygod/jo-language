package parser

import (
	"fmt"

	L "github.com/overlorddamygod/lexer/lexer"
)

type Parser struct {
	lexer *L.Lexer
}

func NewParser(lexer *L.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() {
	// for {
	// 	token, err := p.lexer.NextToken()

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		break
	// 	}

	// 	if token.Type == L.EOF {
	// 		break
	// 	}
	// 	fmt.Printf("%s %s\n", token.Type, token.Literal)
	// }
	e := p.assignment()
	e.Print()
	// e.left()

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

	semicolon, err := p.lexer.NextToken()

	if err != nil {

	}

	if !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		fmt.Println("EXPECTED ;")
		return nil
	}

	return NewBinaryExpression("ASSIGNMENT", identifier, exp)
}

// parse expression
func (p *Parser) expression() Node {
	return p.term()
}

func (p *Parser) factor() Node {
	left := p.primary()

	// if err != nil {

	// }

	for {
		op_token, err := p.lexer.PeekToken(0)

		if err != nil {
			break
		}
		// fmt.Println("HERE")

		// p.lexer.NextToken()

		if op_token.Type == L.OPERATOR && (op_token.Literal == L.SLASH || op_token.Literal == L.ASTERISK) {
			// fmt.Println("OP")
			p.lexer.NextToken()

			op_token.Print()

			right := p.primary()

			left = NewBinaryExpression(string(op_token.Literal), left, right)
		} else {
			// fmt.Println("LOL")
			break
		}
	}
	return left
}

func (p *Parser) primary() Node {
	token, err := p.lexer.NextToken()

	if err != nil {

	}

	if token.Type == L.STRING || token.Type == L.INT || token.Type == L.FLOAT {
		token.Print()
		return NewLiteralValue(string(token.Type), token.Literal)
	}

	if token.Type == L.PUNCTUATION && token.Literal == L.LPAREN {
		token.Print()

		e := p.expression()
		token, err := p.lexer.NextToken()

		if err != nil {

		}
		if token.Type == L.PUNCTUATION && token.Literal == L.RPAREN {
			token.Print()
			return e
		}
	}
	return nil
}

func (p *Parser) term() Node {
	left := p.factor()

	for {
		op_token, err := p.lexer.PeekToken(0)

		if err != nil {
			break
		}

		// p.lexer.NextToken()

		if op_token.Type == L.OPERATOR && (op_token.Literal == L.PLUS || op_token.Literal == L.MINUS) {
			p.lexer.NextToken()
			op_token.Print()

			right := p.factor()
			left = NewBinaryExpression(string(op_token.Literal), left, right)
		} else {
			break
		}
	}
	return left
}
