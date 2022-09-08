package parser

import (
	"errors"
	"fmt"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

type Parser struct {
	lexer *L.Lexer
}

func NewParser(lexer *L.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() ([]node.Node, error) {
	return p.program()
}

func (p *Parser) program() ([]node.Node, error) {
	var declarations []node.Node = make([]node.Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == L.EOF {
			break
		}
		declaration, err := p.declaration()

		if err != nil {
			return declarations, err
		}
		declarations = append(declarations, declaration)
	}
	return declarations, nil
}

func (p *Parser) matchSemicolon(node node.Node) (node.Node, error) {
	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, JoError.New(p.lexer, semicolon, JoError.SyntaxError, "Expected ;")
	}
	return node, nil
}

func (p *Parser) match(type_ L.TokenType, str string) (*L.Token, error) {
	token, err := p.lexer.NextToken()

	if err != nil || (token.Type != type_ || token.Literal != str) {
		return nil, JoError.New(p.lexer, token, JoError.SyntaxError, fmt.Sprintf("Expected ` %s `", str))
	}
	return token, nil
}

func (p *Parser) matchMany(type_ L.TokenType, str ...string) (*L.Token, error) {
	token, err := p.lexer.NextToken()

	if err != nil {
		return nil, errors.New("token not found")
	}

	for _, s := range str {
		if token.Type == type_ && token.Literal == s {
			return token, nil
		}
	}
	return nil, errors.New("no match found")
}
