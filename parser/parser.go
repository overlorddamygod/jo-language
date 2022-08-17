package parser

import (
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
	e := p.expression()
	e.Print()
	// e.left()

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

		if op_token.Type == L.OPERATOR || op_token.Literal == L.SLASH || op_token.Literal == L.ASTERISK {
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

		if op_token.Type == L.OPERATOR || op_token.Literal == L.PLUS || op_token.Literal == L.MINUS {
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
