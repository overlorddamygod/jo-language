package parser

import (
	"fmt"

	L "github.com/overlorddamygod/jo/pkg/lexer"
)

type Parser struct {
	lexer *L.Lexer
}

func NewParser(lexer *L.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() ([]Node, error) {
	return p.program()
}

func (p *Parser) program() ([]Node, error) {
	var declarations []Node = make([]Node, 0)
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

func (p *Parser) declaration() (Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "fn":
		return p.functionDecl()
	case "struct":
		return p.structDecl()
	case "let":
		decl, err := p.vardecl()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(decl)
	}
	return p.statement()
	// return nil, L.NewJoError(p.lexer, first, fmt.Sprintf("Unknown declaration ` %s `", first.Literal))
}

func (p *Parser) structDecl() (Node, error) {
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

	var methods []FunctionDeclStatement

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

	return NewStructDeclStatement(identifier, methods), nil
}

func (p *Parser) vardecl() (Node, error) {
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

	if err != nil {
		return nil, err
	}

	return NewVarDeclStatement(identifier, expression), nil
}

func (p *Parser) statements() ([]Node, error) {
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
func (p *Parser) declarations() ([]Node, error) {
	var declarations []Node = make([]Node, 0)
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

func (p *Parser) statement() (Node, error) {
	first, _ := p.lexer.PeekToken(0)

	switch first.Literal {
	case "if":
		return p.ifElse()
	case "for":
		return p.For()
	case "return":
		ret, err := p._return()

		if err != nil {
			return nil, err
		}
		return p.matchSemicolon(ret)
	// TODO: ADD BLOCK STATEMENT
	// case "{":
	// return p.block()
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

func (p *Parser) _break() (Node, error) {
	_, err := p.match(L.KEYWORD, "break")

	if err != nil {
		return nil, err
	}

	return NewBreakStatement(), nil
}

func (p *Parser) _continue() (Node, error) {
	_, err := p.match(L.KEYWORD, "continue")

	if err != nil {
		return nil, err
	}

	return NewContinueStatement(), nil
}

func (p *Parser) _return() (Node, error) {
	_, err := p.match(L.KEYWORD, "return")

	if err != nil {
		return nil, err
	}

	semicolon, _ := p.lexer.PeekToken(0)

	if semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON {
		return NewReturnStatement(nil), nil
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	return NewReturnStatement(exp), nil
}

func (p *Parser) matchSemicolon(node Node) (Node, error) {
	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, L.NewJoError(p.lexer, semicolon, "Expected ;")
	}
	return node, nil
}

func (p *Parser) functionDecl() (*FunctionDeclStatement, error) {
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

		return NewFunctionDeclStatement(identifier, []Node{}, block), nil
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
	return NewFunctionDeclStatement(identifier, parameters, block), nil
}

func (p *Parser) match(type_ L.TokenType, str string) (*L.Token, error) {
	token, err := p.lexer.NextToken()

	if err != nil || (token.Type != type_ || token.Literal != str) {
		return nil, L.NewJoError(p.lexer, token, fmt.Sprintf("Expecred ` %s `", str))
	}
	return token, nil
}

func (p *Parser) condition() (Node, error) {
	_, err := p.match(L.PUNCTUATION, L.LPAREN)

	if err != nil {
		return nil, err
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}

	_, err = p.match(L.PUNCTUATION, L.RPAREN)

	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) ifElse() (Node, error) {
	_, err := p.match(L.KEYWORD, "if")

	if err != nil {
		return nil, err
	}

	exp, err := p.condition()

	if err != nil {
		return nil, err
	}

	// fmt.Println("HERE")

	ifBlock, err := p.block()

	if err != nil {
		return nil, err
	}

	var ifs []*ConditionBlock = make([]*ConditionBlock, 0)
	ifs = append(ifs, NewConditionBlock(exp, ifBlock))

	token, _ := p.lexer.PeekToken(0)

	if token.Literal == "elif" {

		for {
			p.lexer.NextToken()

			exp, err := p.expression()

			if err != nil {
				return nil, err
			}

			block, err := p.block()

			if err != nil {
				return nil, err
			}

			ifs = append(ifs, NewConditionBlock(exp, block))

			token, _ = p.lexer.PeekToken(0)

			if token.Literal != "elif" {
				break
			}
		}
	}

	ifStatement := NewIfStatement(ifs)

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
		return nil, L.NewJoError(p.lexer, identifier, "Expected for")
	}

	leftParenthesis, err := p.lexer.NextToken()

	if err != nil || !(leftParenthesis.Type == L.PUNCTUATION && leftParenthesis.Literal == L.LPAREN) {
		return nil, L.NewJoError(p.lexer, leftParenthesis, "Expected (")
	}

	vardecl, err := p.vardecl()

	if err != nil {
		return nil, err
	}

	semicolon, err := p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, L.NewJoError(p.lexer, semicolon, "Expected ;")
	}

	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	// condition.Print()

	semicolon, err = p.lexer.NextToken()

	if err != nil || !(semicolon.Type == L.PUNCTUATION && semicolon.Literal == L.SEMICOLON) {
		return nil, L.NewJoError(p.lexer, semicolon, "Expected ;")
	}

	exp, err := p.expression()

	if err != nil {
		return nil, err
	}
	// exp.Print()

	rightParenthesis, err := p.lexer.NextToken()

	if err != nil || !(rightParenthesis.Type == L.PUNCTUATION && rightParenthesis.Literal == L.RPAREN) {
		return nil, L.NewJoError(p.lexer, semicolon, "Expected )")
	}
	// fmt.Println("HERE")

	block, err := p.block()

	if err != nil {
		return nil, err
	}

	return NewForStatement(vardecl, condition, exp, block), nil
}

func (p *Parser) block() (*Block, error) {
	leftCurly, err := p.lexer.NextToken()

	// fmt.Println("LEFTCURLY", leftCurly)
	if err != nil || !(leftCurly.Type == L.PUNCTUATION && leftCurly.Literal == L.LBRACE) {
		return nil, L.NewJoError(p.lexer, leftCurly, "Expected {")
	}

	// fmt.Println("BLOCK")

	block, err := p.declarations()

	if err != nil {
		return nil, err
	}

	rightCurly, err := p.lexer.NextToken()

	if err != nil || !(rightCurly.Type == L.PUNCTUATION && rightCurly.Literal == L.RBRACE) {
		return nil, L.NewJoError(p.lexer, rightCurly, "Expected }")
	}

	return NewBlock(block), nil
}

func (p *Parser) identifier() (Node, error) {
	identifier, err := p.lexer.NextToken()

	if err != nil || identifier.Type == L.EOF || identifier.Type != L.IDENTIFIER {
		return nil, L.NewJoError(p.lexer, identifier, "Expected identifier")
	}

	if L.IsKeyword(identifier.Literal) {
		return nil, L.NewJoError(p.lexer, identifier, "Variable name cannot be a keyword")
	}

	return NewIdentifier(identifier.Literal, identifier), nil
}

// parse expression
func (p *Parser) expression() (Node, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Node, error) {
	identi, _ := p.lexer.PeekToken(0)

	if identi.Type == L.IDENTIFIER {
		equal, _ := p.lexer.PeekToken(1)

		if equal.Type == L.OPERATOR && equal.Literal == L.ASSIGN {
			identifier, err := p.identifier()

			if err != nil {
				return nil, err
			}

			equals, err := p.lexer.NextToken()

			if err != nil || !(equals.Type == L.OPERATOR && equals.Literal == "=") {
				return nil, L.NewJoError(p.lexer, equals, "Expected =")
			}

			exp, err := p.expression()

			if err != nil {
				return nil, err
			}

			return NewAssignmentStatement(identifier, exp), nil
		}
	}
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
	return p.binary(p.unary, func(t *L.Token) bool {
		return t.Type == L.OPERATOR && (t.Literal == L.SLASH || t.Literal == L.ASTERISK || t.Literal == L.PERCENT)
	})
}

func (p *Parser) unary() (Node, error) {
	token, _ := p.lexer.PeekToken(0)

	if token.Type == L.OPERATOR && (token.Literal == L.BANG || token.Literal == L.UNARY_PLUS || token.Literal == L.UNARY_MINUS) {
		p.lexer.NextToken()
		unary_, _ := p.unary()
		return NewUnaryExpression(string(token.Literal), unary_, token), nil
	}

	return p.call()
}

func (p *Parser) call() (Node, error) {
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

				expr = NewFunctionCall(expr, []Node{})
				continue
			}

			arguments, err := p.arguments()

			if err != nil {
				return nil, err
			}

			// fmt.Println(arguments)

			rightParen, err := p.lexer.PeekToken(0)

			if err != nil || !(rightParen.Type == L.PUNCTUATION && rightParen.Literal == L.RPAREN) {
				return nil, L.NewJoError(p.lexer, rightParen, "Expected )")
			}
			p.lexer.NextToken()
			expr = NewFunctionCall(expr, arguments)
		} else if leftParenOrDot.Type == L.OPERATOR && leftParenOrDot.Literal == L.FULL_STOP {
			p.lexer.NextToken()

			iden, err := p.identifier()

			if err != nil {
				return nil, err
			}

			expr = NewGetExpr(iden, expr)
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) primary() (Node, error) {
	token, err := p.lexer.NextToken()

	// fmt.Println("BOOO", token)

	if err != nil {
		return nil, L.NewJoError(p.lexer, token, "Expected value")
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
		return NewIdentifier(token.Literal, token), nil
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
			return nil, L.NewJoError(p.lexer, token, "Expected )")

		}
		return e, nil
	}
	return nil, L.NewJoError(p.lexer, token, "unexpected value")
}
