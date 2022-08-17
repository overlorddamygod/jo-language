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

func (l *Lexer) advanceWithToken(token Token) {
	l.tokens = append(l.tokens, token)
	l.pos += 1
	l.col += 1
}

func (l *Lexer) appendToken(token Token) {
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) getStringLiteral() (string, error) {
	startPos := l.pos

	current, err := l.current()

	if err != nil {
		return "", err
	}
	prev := current
	l.advance()

	for {
		current, err := l.current()

		if err != nil {
			return "", err
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

	return l.source[startPos:l.pos], nil
}

func (l *Lexer) getIdentifier() (string, error) {
	startPos := l.pos

	current, err := l.current()

	if err != nil {
		return "", err
	}
	for IsAlphaNumericWithUnderscore(current) {
		l.advance()

		current, err = l.current()

		if err != nil {
			return "", err
		}
	}

	return l.source[startPos:l.pos], nil
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

func (l *Lexer) getNumberLiteral() (string, bool, error) {
	startPos := l.pos
	hasDecimal := false

	for {
		current, err := l.current()

		if err != nil {
			return l.source[startPos:l.pos], hasDecimal, nil
		}

		if !(IsDigit(current) || current == ".") {
			return l.source[startPos:l.pos], hasDecimal, nil
		}

		if current == "." {
			if hasDecimal {
				return "", hasDecimal, errors.New("failed parsing as number")
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
		} else {
			l.col += 1
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
		fmt.Println("BEFORE SKIPPING", l.pos)

		err := l.skipWhiteSpace()
		fmt.Println("AFTER SKIPPING", l.pos, l.size)

		if err != nil {
			break
		}

		currentChar, err := l.current()
		fmt.Println(currentChar, l.pos)

		if err != nil {
			return l.tokens, err
		}

		if IsLetter(currentChar) || currentChar == "_" {
			identifier, err := l.getIdentifier()
			if err != nil {
				return l.tokens, err
			}

			var tokenType TokenType = IDENTIFIER

			if IsKeyword(identifier) {
				tokenType = KEYWORD
			}

			l.appendToken(NewToken(tokenType, identifier))
			continue
		}

		if IsDigit(currentChar) {
			num, hasDecimal, err := l.getNumberLiteral()
			fmt.Println("NUM", num, l.pos)
			if err != nil {
				return l.tokens, err
			}
			var tokenType TokenType = INT

			if hasDecimal {
				tokenType = FLOAT
			}
			l.appendToken(NewToken(tokenType, num))
			continue
		}

		switch currentChar {
		case PLUS, MINUS, ASTERISK, SLASH, PERCENT:
			fmt.Println("OPERATOR", currentChar, l.pos)

			l.advanceWithToken(NewToken(OPERATOR, currentChar))
			fmt.Println("OPERATOR", currentChar, l.pos)

			continue

		case ASSIGN, EQ:
			peek, err := l.peek(1)
			if err != nil {
				break
			}
			if peek == "=" {
				l.advanceWithToken(NewToken(OPERATOR, EQ))
			} else {
				l.appendToken(NewToken(OPERATOR, ASSIGN))
			}
			l.advance()
			continue

		case BANG, NOT_EQ:
			peek, err := l.peek(1)

			if err != nil {
				break
			}

			if peek == "=" {
				l.advanceWithToken(NewToken(OPERATOR, NOT_EQ))
			} else {
				l.appendToken(NewToken(OPERATOR, BANG))
			}
			l.advance()
			continue
		case LPAREN, RPAREN, LBRACE, RBRACE, COMMA, SEMICOLON, COLON, LT, GT:
			l.advanceWithToken(NewToken(PUNCTUATION, currentChar))
			continue
		case DOUBLE_QUOTE:
			strLiteral, err := l.getStringLiteral()

			if err != nil {
				return l.tokens, err
			}

			l.appendToken(NewToken(STRING, strLiteral))
			continue
		default:
			return l.tokens, fmt.Errorf("-- Line: %d Col: %d Pos %d: Illegal Character `%s`", l.line, l.col, l.pos, currentChar)
		}
		fmt.Println("HERE")
		l.advance()
	}
	l.appendToken(NewToken(EOF, "EOF"))
	l.token_size = len(l.tokens)
	return l.tokens, nil
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