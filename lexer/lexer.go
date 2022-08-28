package lexer

import (
	"errors"
	"fmt"
)

type Lexer struct {
	source     string
	start      int
	col        int
	pos        int
	size       int
	line       int
	tokens     []Token
	token_pos  int
	token_size int
}

func NewLexer(source string) *Lexer {
	var tokens []Token = make([]Token, 0)
	return &Lexer{source: source, start: 0, pos: 0, size: len(source), line: 1, col: 0, tokens: tokens, token_pos: 0, token_size: 0}
}

func (l *Lexer) peek(offset int) (string, error) {
	if l.pos+offset >= l.size {
		return "", errors.New("EOF")
	}
	return string(l.source[l.pos+offset]), nil
}

func (l *Lexer) current() (string, error) {
	return l.peek(0)
}

func (l *Lexer) advance() {
	l.pos += 1
	l.col += 1
}

func (l *Lexer) advanceWithToken(token *Token) {
	l.tokens = append(l.tokens, *token)
	l.pos += 1
	l.col += 1
	// token.End(l.col)
}

func (l *Lexer) appendToken(token *Token) {
	l.tokens = append(l.tokens, *token)
}

func (l *Lexer) addToken(tokenType TokenType, literal string) {
	token := NewToken(tokenType, literal).Start(l.start + 1).End(l.pos)

	l.tokens = append(l.tokens, *token)
}

func (l *Lexer) getToken(tokenType TokenType, literal string) *Token {
	token := NewToken(tokenType, literal).Line(l.line).Start(l.start + 1).End(l.col)
	return token
}
func (l *Lexer) getStringLiteral() (*Token, error) {
	startPos := l.pos
	current, err := l.current()

	if err != nil {
		return nil, err
	}
	prev := current
	l.advance()

	for {
		current, err := l.current()

		if err != nil || current == "\n" {
			return l.getToken(STRING, l.source[startPos+1:l.pos-1]), errors.New("expected '\"'")
		}

		if current == "\"" {
			if prev != "\\" {
				l.advance()
				break
			}
		}

		prev = current

		l.advance()
	}
	return l.getToken(STRING, l.source[startPos+1:l.pos-1]), nil
}

func (l *Lexer) getIdentifier() (*Token, error) {
	startPos := l.pos

	current, err := l.current()

	if err != nil {
		return nil, err
	}
	for IsAlphaNumericWithUnderscore(current) {
		l.advance()

		current, err = l.current()

		if err != nil {
			return l.getToken(IDENTIFIER, l.source[startPos:l.pos]), err
		}
	}

	return l.getToken(IDENTIFIER, l.source[startPos:l.pos]), nil
}

// func (l *Lexer) getNumberLiteral() (string, bool, error) {
// 	startPos := l.pos

// 	current, err := l.current()

// 	if err != nil {
// 		return "", false, err
// 	}

// 	hasDecimal := false

// 	for IsDigit(current) || current == "." {

// 		if current == "." {
// 			if hasDecimal {
// 				return "", hasDecimal, errors.New("failed parsing as number")
// 			}
// 			hasDecimal = true
// 		}

// 		l.advance()

// 		current, err = l.current()

// 		if err != nil {
// 			return "", hasDecimal, err
// 		}
// 	}

// 	return l.source[startPos:l.pos], hasDecimal, nil
// }

func (l *Lexer) getNumberLiteral() (*Token, error) {
	startPos := l.pos
	hasDecimal := false
	var tokenType TokenType = INT
	for {
		current, err := l.current()

		if err != nil {
			if hasDecimal {
				tokenType = FLOAT
			}
			return l.getToken(tokenType, l.source[startPos:l.pos]), nil
		}

		if !(IsDigit(current) || current == ".") {
			if hasDecimal {
				tokenType = FLOAT
			}
			return l.getToken(tokenType, l.source[startPos:l.pos]), nil

			// return l.source[startPosPos:l.pos], hasDecimal, nil
		}

		if current == "." {
			if hasDecimal {
				return l.getToken(tokenType, l.source[startPos:l.pos]), errors.New("failed parsing as number")
			}
			hasDecimal = true
		}

		l.advance()
	}
}

func (l *Lexer) skipWhiteSpace() error {
	current, err := l.current()

	if err != nil {
		return err
	}
	for current == " " || current == "\n" || current == "\t" || current == "\r" {
		l.advance()
		if current == "\n" {
			l.line += 1
			l.col = 0
		}

		current, err = l.current()

		if err != nil {
			return err
		}
	}
	return nil
}

func IsLetter(s string) bool {
	if (s < "a" || s > "z") && (s < "A" || s > "Z") {
		return false
	}
	return true
}

func IsDigit(s string) bool {
	return s >= "0" && s <= "9"
}

func IsAlphaNumericWithUnderscore(s string) bool {
	return IsLetter(s) || IsDigit(s) || s == "_"
}

func (l *Lexer) Lex() ([]Token, error) {
	for l.pos < l.size {
		err := l.skipWhiteSpace()

		if err != nil {
			break
		}

		l.start = l.col

		currentChar, err := l.current()

		if err != nil {
			return l.tokens, err
		}

		if IsLetter(currentChar) || currentChar == "_" {
			identifier, _ := l.getIdentifier()

			if IsKeyword(identifier.Literal) {
				identifier.SetType(KEYWORD)
			}

			l.appendToken(identifier)
			continue
		}

		if IsDigit(currentChar) {
			num, err := l.getNumberLiteral()
			if err != nil {
				return l.tokens, NewJoError(l, num, err.Error())
			}
			l.appendToken(num)
			continue
		}

		switch currentChar {
		case ASTERISK, PERCENT:
			l.advance()
			l.appendToken(l.getToken(OPERATOR, currentChar).Line(l.line).Start(l.col + 1).End(l.col + 1))
			continue

		case SLASH:
			peek, err := l.peek(1)
			if err != nil {
				break
			}
			if peek == SLASH {
				for currentChar != "\n" {
					l.advance()
					if l.pos >= l.size {
						break
					}

					currentChar, _ = l.current()
				}
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, SLASH))
			}
			continue
		case ASSIGN:
			peek, err := l.peek(1)
			if err != nil {
				break
			}
			if peek == "=" {
				l.advance()
				l.advance()
				// l.appendToken(l.getToken(OPERATOR, EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
				l.appendToken(l.getToken(OPERATOR, EQ))
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, ASSIGN))
			}
			continue
		case LT, GT:
			peek, err := l.peek(1)
			if err != nil {
				// break
			}
			if peek == "=" {
				l.advance()
				l.advance()
				literal := LT_EQ
				if currentChar == GT {
					literal = GT_EQ
				}
				// l.appendToken(l.getToken(OPERATOR, EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
				l.appendToken(l.getToken(OPERATOR, literal))
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, currentChar))
			}
			continue
		case PIPE, AMPERSAND:
			peek, err := l.peek(1)
			if err != nil {
				// break
			}
			if peek == currentChar {
				l.advance()
				l.advance()

				literal := OR

				if currentChar == AMPERSAND {
					literal = AND
				}
				// l.appendToken(l.getToken(OPERATOR, EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
				l.appendToken(l.getToken(OPERATOR, literal))
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, currentChar))
			}
			continue
		case PLUS, MINUS:
			peek, err := l.peek(1)
			if err != nil {
				// break
			}
			if peek == currentChar {
				l.advance()
				l.advance()

				literal := UNARY_PLUS

				if currentChar == MINUS {
					literal = UNARY_MINUS
				}
				// l.appendToken(l.getToken(OPERATOR, EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
				l.appendToken(l.getToken(OPERATOR, literal))
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, currentChar))
			}
			continue

		case BANG, NOT_EQ:
			peek, err := l.peek(1)

			if err != nil {
				break
			}

			if peek == "=" {
				l.advance()
				l.advance()
				l.appendToken(l.getToken(OPERATOR, NOT_EQ))
			} else {
				l.advance()
				l.appendToken(l.getToken(OPERATOR, BANG))
			}
			continue
		case LPAREN, RPAREN, LBRACE, RBRACE, COMMA, SEMICOLON, COLON:
			l.advance()
			l.appendToken(l.getToken(PUNCTUATION, currentChar))

			continue
		case DOUBLE_QUOTE:
			strLiteralToken, err := l.getStringLiteral()

			if err != nil {
				return l.tokens, NewJoError(l, strLiteralToken, fmt.Sprintf("expected ` %s `", currentChar))
			}

			l.appendToken(strLiteralToken)
			continue
		default:
			l.advance()
			token := l.getToken(ILLEGAL, currentChar)
			return l.tokens, NewJoError(l, token, fmt.Sprintf("Illegal Character `%s`", currentChar))
		}
		// fmt.Println("HERE")
		l.advance()
	}

	l.start = l.col
	l.advance()
	l.appendToken(l.getToken(EOF, "EOF"))
	l.token_size = len(l.tokens)
	return l.tokens, nil
}

func (l *Lexer) GetLine(line int) string {
	src := ""
	pos := 0
	lineNo := 1
	for {
		if lineNo == line {
			for {
				if pos == l.size-1 || string(l.source[pos]) == "\n" {
					src += string(l.source[pos])
					return fmt.Sprintf(" %d | %s", line, src)
				}
				src += string(l.source[pos])
				pos++
			}
		}
		// loop until we find a newline
		if string(l.source[pos]) == "\n" {
			lineNo++
		}
		pos++

		if pos >= l.size {
			break
		}
	}

	return src
}

func (l *Lexer) NextToken() (*Token, error) {
	if l.token_pos >= l.token_size {
		return nil, errors.New("EOT")
	}
	l.token_pos += 1
	return &l.tokens[l.token_pos-1], nil
}

func (l *Lexer) PeekToken(offset int) (*Token, error) {
	if l.token_pos+offset >= l.token_size {
		return nil, errors.New("EOT")
	}
	return &l.tokens[l.token_pos+offset], nil
}
