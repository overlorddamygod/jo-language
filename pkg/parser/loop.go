package parser

import (
	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) For() (node.Node, error) {
	// fmt.Println("ERE")
	identifier, err := p.lexer.NextToken()

	if err != nil || identifier.Literal != "for" {
		return nil, JoError.New(p.lexer, identifier, JoError.SyntaxError, "Expected for")
	}

	leftParenthesis, err := p.lexer.NextToken()
	// fmt.Println(leftParenthesis)
	if err != nil || !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		return nil, JoError.New(p.lexer, leftParenthesis, JoError.SyntaxError, "Expected (")
	}

	vardecl, err := p.vardecl()

	if err != nil {
		return nil, err
	}

	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, JoError.New(p.lexer, semicolon, JoError.SyntaxError, "Expected ;")
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	// condition.Print()

	semicolon, err = p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, JoError.New(p.lexer, semicolon, JoError.SyntaxError, "Expected ;")
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}
	// exp.Print()

	rightParenthesis, err := p.lexer.NextToken()

	if err != nil || !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		return nil, JoError.New(p.lexer, semicolon, JoError.SyntaxError, "Expected )")
	}
	// fmt.Println("HERE")

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return node.NewForStatement(vardecl, condition, exp, block), nil
}

func (p *Parser) While() (node.Node, error) {
	if _, err := p.match(L.KEYWORD, "while"); err != nil {
		return nil, err
	}

	if _, err := p.match(L.PUNCTUATION, L.LPAREN); err != nil {
		return nil, err
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.match(L.PUNCTUATION, L.RPAREN); err != nil {
		return nil, err
	}

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return node.NewWhileStatement(condition, block), nil
}
