package parser

import (
	"errors"

	"github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

func (p *Parser) statement() (node.Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "import":
		return p.importStatement()
	case "export":
		return p.exportStatement()
	case "if":
		return p.ifElse()
	case "while":
		return p.While()
	case "for":
		return p.For()
	case "return":
		ret, err := p._return()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(ret)
	case "try":
		return p.tryCatch()
	case "throw":
		throw, err := p.throw()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(throw)
	case "{":
		return p.block()
	case "break":
		br, err := p._break()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(br)

	case "continue":
		con, err := p._continue()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(con)
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	return p.matchSemicolon(exp)
}

func (p *Parser) importStatement() (node.Node, error) {
	_, err := p.match(lexer.KEYWORD, "import")
	if err != nil {
		return nil, err
	}

	fileName, err := p.lexer.NextToken()

	if err != nil {
		return nil, errors.New("Expected import file name")
	}

	if fileName.Type != lexer.STRING {
		return nil, errors.New("Expected import file name")
	}

	return node.NewImport(fileName), nil
}

func (p *Parser) exportStatement() (node.Node, error) {
	_, err := p.match(lexer.KEYWORD, "export")

	if err != nil {
		return nil, err
	}

	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if err := p.peekMatch(0, lexer.PUNCTUATION, lexer.SEMICOLON); err == nil {
		p.lexer.NextToken()
	}

	return node.NewExport(expr), nil
}
