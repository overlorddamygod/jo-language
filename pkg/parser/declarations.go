package parser

import (
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) declarations() ([]node.Node, error) {
	var declarations []node.Node = make([]node.Node, 0)
	for {
		token, _ := p.lexer.PeekToken(0)
		// fmt.Println("SAD", token.Type)
		if token.Type == L.PUNCTUATION && token.Literal == L.RBRACE {
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

func (p *Parser) declaration() (node.Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "fn":
		return p.functionDecl()
	case "struct":
		return p.structDecl()
	case "switch":
		return p.switchDecl()
	case "let":
		decl, err := p.vardecl()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(decl)
	}
	return p.statement()
	// return nil, JoError.New(p.lexer, first, JoError.SyntaxError, fmt.Sprintf("Unknown declaration ` %s `", first.Literal))
}

func (p *Parser) structDecl() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "struct")

	if err != nil {
		return nil, err
	}

	identifier, err := p.identifier()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.LBRACE)

	if err != nil {
		return nil, err
	}

	var methods []node.FunctionDeclStatement

	for {
		token, err := p.lexer.PeekToken(0)

		if err != nil {
			break
		}

		if token.Literal == "fn" {
			method, err := p.functionDecl()
			// method.Print()
			if err != nil {
				return nil, err
			}
			methods = append(methods, *method)
		} else {
			break
		}
	}

	_, err = p.match(L.PUNCTUATION, L.RBRACE)

	if err != nil {
		return nil, err
	}

	return node.NewStructDeclStatement(identifier, methods), nil
}

func (p *Parser) vardecl() (node.Node, error) {
	_, err := p.match(L.KEYWORD, "let")
	if err != nil {
		return nil, err
	}

	identifier, err := p.identifier()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.OPERATOR, L.ASSIGN)

	if err != nil {
		return nil, err
	}

	expression, err := p.expression()

	// fmt.Println("SAD", expression, err)
	if err != nil {
		return nil, err
	}

	return node.NewVarDeclStatement(identifier, expression), nil
}
