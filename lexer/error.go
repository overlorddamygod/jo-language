package lexer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type JoError struct {
	token    *Token
	errorMsg error
}

func (m *JoError) Error() string {
	return m.errorMsg.Error()
}

func NewJoError(l *Lexer, token *Token, msg string) *JoError {
	line, err := l.GetLine(token.line)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &JoError{
		token:    token,
		errorMsg: errors.New(MarkError(line, token.line, token.start, token.end, msg)),
	}
}

func MarkError(line string, lineNo int, start int, end int, msg string) string {
	lineNoLength := len(strconv.Itoa(lineNo))

	gap := strings.Repeat(" ", 3+lineNoLength+start)
	marker := strings.Repeat("^", end-start+1)

	return fmt.Sprintf("%d | %s   \n%s%s   \n-- Line: %d Col: %d : %s\n", lineNo, line, gap, marker, lineNo, start, msg)
}
