package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Lexer struct {
	source     string
	src        string
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
	return &Lexer{source: source, src: source, start: 0, pos: 0, size: len(source), line: 1, col: 0, tokens: tokens, token_pos: 0, token_size: 0}
}

func (l *Lexer) Lex() ([]Token, *Token, error) {
	for {
		token, err := l.Next()
		if err != nil {
			if token.Type != EOF {
				return l.tokens, token, err
			}
		}
		l.tokens = append(l.tokens, *token)

		if token.Type == EOF {
			break
		}
	}
	l.token_size = len(l.tokens)
	return l.tokens, nil, nil
}

func (l *Lexer) GetLine(line int) (string, error) {
	split := strings.Split(l.source, "\n")

	if line > len(split) {
		return "", errors.New("line greater than the total lines")
	}
	if line <= 0 {
		return "", errors.New("line less or equal to zero")
	}

	return split[line-1], nil
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

func (l *Lexer) SetTokenPos(pos int) {
	l.token_pos = pos
}

func (l *Lexer) GetTokenPos() int {
	return l.token_pos
}

func (l *Lexer) hasString() bool {
	return len(l.src) > 0
}

func (l *Lexer) Next() (*Token, error) {
	start := l.col + 1

	if !l.hasString() {
		return NewToken(EOF, "", l.line, start, start), errors.New("EOF")
	}

	for _, regexs := range RegexTokenMap {
		regex, tokenType := regexs[0], regexs[1]

		re := regexp.MustCompile(regex)
		matchedString := re.FindString(l.src)

		lenMatch := len(matchedString)
		if lenMatch > 0 {
			l.src = l.src[lenMatch:]

			if tokenType == WHITESPACE || tokenType == COMMENT {
				l.col += lenMatch
				return l.Next()
			}
			if tokenType == NEWLINE {
				l.col = 0
				l.line += 1
				return l.Next()
			}

			l.col += lenMatch

			if tokenType == IDENTIFIER && IsKeyword(matchedString) {
				return NewToken(KEYWORD, matchedString, l.line, start, l.col), nil
			}
			if tokenType == STRING {
				return NewToken(STRING, matchedString[1:len(matchedString)-1], l.line, start, l.col), nil
			}
			return NewToken(TokenType(tokenType), matchedString, l.line, start, l.col), nil
		}
	}

	return NewToken(ILLEGAL, l.src[0:1], l.line, start, start), fmt.Errorf("%s character `%s`", ILLEGAL, l.src[0:1])
}
