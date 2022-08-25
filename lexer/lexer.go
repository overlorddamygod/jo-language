package lexer

import (
	"errors"
	"fmt"
)

type Lexer struct {
	source     string
	pos        int
	size       int
	line       int
	col        int
	tokens     []Token
	token_pos  int
	token_size int
}

func NewLexer(source string) *Lexer {
	var tokens []Token = make([]Token, 0)
	return &Lexer{source: source, pos: 0, size: len(source), line: 1, col: 0, tokens: tokens, token_pos: 0, token_size: 0}
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

		if err != nil {
			return NewToken(STRING, l.source[startPos:l.pos]).Line(l.line).Start(startPos + 1).End(l.pos), errors.New("expected '\"'")
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
	return NewToken(STRING, l.source[startPos:l.pos]).Line(l.line).Start(startPos + 1).End(l.pos), nil
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
			return nil, err
		}
	}

	id := l.source[startPos:l.pos]

	return NewToken(IDENTIFIER, id).Line(l.line).Start(startPos + 1).End(l.pos), nil
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
			return NewToken(tokenType, l.source[startPos:l.pos]).Line(l.line).Start(startPos + 1).End(l.pos), nil
		}

		if !(IsDigit(current) || current == ".") {
			if hasDecimal {
				tokenType = FLOAT
			}
			return NewToken(tokenType, l.source[startPos:l.pos]).Line(l.line).Start(startPos + 1).End(l.pos), nil

			// return l.source[startPos:l.pos], hasDecimal, nil
		}

		if current == "." {
			if hasDecimal {
				return NewToken(tokenType, l.source[startPos:l.pos]).Line(l.line).Start(startPos + 1).End(l.pos), errors.New("failed parsing as number")
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
		if current == "\n" {
			l.line += 1
			l.col = 0
		}
		l.advance()

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

		currentChar, err := l.current()

		if err != nil {
			return l.tokens, err
		}

		if IsLetter(currentChar) || currentChar == "_" {
			identifier, err := l.getIdentifier()
			if err != nil {
				return l.tokens, err
			}

			if IsKeyword(identifier.Literal) {
				identifier.SetType(KEYWORD)
			}

			l.appendToken(identifier)
			continue
		}

		if IsDigit(currentChar) {
			num, err := l.getNumberLiteral()
			if err != nil {
				return l.tokens, fmt.Errorf("%s", MarkError(l.GetLine(num.line), num.line, num.start, l.col+1, err.Error()))
			}
			l.appendToken(num)
			continue
		}

		switch currentChar {
		case PLUS, MINUS, ASTERISK, SLASH, PERCENT:
			l.advanceWithToken(NewToken(OPERATOR, currentChar).Line(l.line).Start(l.col + 1).End(l.col + 1))
			continue

		case ASSIGN, EQ:
			peek, err := l.peek(1)
			if err != nil {
				break
			}
			if peek == "=" {
				l.advanceWithToken(NewToken(OPERATOR, EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
			} else {
				l.appendToken(NewToken(OPERATOR, ASSIGN).Line(l.line).Start(l.col + 1).End(l.col + 1))
			}
			l.advance()
			continue

		case BANG, NOT_EQ:
			peek, err := l.peek(1)

			if err != nil {
				break
			}

			if peek == "=" {
				l.advanceWithToken(NewToken(OPERATOR, NOT_EQ).Line(l.line).Start(l.col + 1).End(l.col + 2))
			} else {
				l.appendToken(NewToken(OPERATOR, BANG).Line(l.line).Start(l.col + 1).End(l.col + 1))
			}
			l.advance()
			continue
		case LPAREN, RPAREN, LBRACE, RBRACE, COMMA, SEMICOLON, COLON, LT, GT:
			l.advanceWithToken(NewToken(PUNCTUATION, currentChar).Line(l.line).Start(l.col + 1).End(l.col + 1))
			continue
		case DOUBLE_QUOTE:
			strLiteralToken, err := l.getStringLiteral()

			if err != nil {
				line := l.GetLine(strLiteralToken.line)
				fmt.Println(strLiteralToken.end, len(line))
				end := strLiteralToken.end
				if strLiteralToken.end >= len(line) {
					end = len(line)
				}
				fmt.Println(strLiteralToken.start, end)
				return l.tokens, fmt.Errorf("%s", MarkError(line, strLiteralToken.line, strLiteralToken.start, end, err.Error()))
			}

			l.appendToken(strLiteralToken)
			continue
		default:
			// l.GetLine(l.line)
			err := fmt.Sprintf("Illegal Character `%s`", currentChar)
			return l.tokens, fmt.Errorf("%s", MarkError(l.GetLine(l.line), l.line, l.col+1, l.col+1, err))
		}
		l.advance()
	}
	l.appendToken(NewToken(EOF, "EOF").Line(l.line).Start(l.col + 1).End(l.col + 1))
	l.token_size = len(l.tokens)
	return l.tokens, nil
}

func MarkError(line string, lineNo int, start int, end int, msg string) string {
	strlen := len(line)
	formatStr := "%s\n"

	for i := 0; i < strlen; i++ {
		if string(line[i]) == "|" {
			formatStr += " "
			break
		}
		formatStr += " "
	}

	for i := 0; i < start; i++ {
		formatStr += " "
	}
	// fmt.Println(start, end)
	for i := start; i <= end; i++ {
		formatStr += "^"
	}

	formatStr += "\n-- Line: %d Col: %d : %s\n"

	return fmt.Sprintf(formatStr, line, lineNo, start, msg)
}

func (l *Lexer) GetLine(line int) string {

	// fmt.Printf("Error on line %d\n", line)

	// for i := 0; i < line; i++ {
	// 	fmt.Println(l.source[i])
	// }
	src := ""
	pos := 0
	lineNo := 1
	for {
		// fmt.Println("HERE", lineNo, line)

		if lineNo == line {
			for {
				if string(l.source[pos]) == "\n" || pos >= l.size {
					return fmt.Sprintf(" %d | %s", line, src)
				}
				src += string(l.source[pos])
				// fmt.Println(string(l.source[pos]))
				pos++
			}
		}
		// loop until we find a newline
		if string(l.source[pos]) == "\n" {
			// fmt.Println("line")

			lineNo++
		}

		// src += l.source[l.pos]
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
