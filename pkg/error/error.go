package joerror

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

type JoError struct {
	_type    JoErrorType
	token    *lexer.Token
	errorMsg error
}

type JoErrorType string

var (
	DefaultError   JoErrorType = "Error"
	LexicalError   JoErrorType = "LexicalError"
	SyntaxError    JoErrorType = "SyntaxError"
	ReferenceError JoErrorType = "ReferenceError"
)

func (m *JoError) Error() string {
	return m.errorMsg.Error()
}

func New(l *lexer.Lexer, token *lexer.Token, _type JoErrorType, msg string) *JoError {
	var err error = fmt.Errorf("[%s] %s", _type, msg)

	joErr := &JoError{
		token:    token,
		_type:    _type,
		errorMsg: err,
	}
	if token != nil {
		line, err := l.GetLine(token.GetLine())
		if err != nil {
			fmt.Println(err.Error())
		}
		err = errors.New(joErr.MarkError(line, token.GetLine(), token.GetStart(), token.GetEnd(), msg))
		joErr.errorMsg = err
	}
	return joErr
}

func (e *JoError) MarkError(line string, lineNo int, start int, end int, msg string) string {
	lineNoLength := len(strconv.Itoa(lineNo))

	gap := strings.Repeat(" ", 2+lineNoLength+start)
	marker := strings.Repeat("^", end-start+1)

	return fmt.Sprintf("%d | %s   \n%s%s   \n-- [ %s ] Line: %d Col: %d : %s\n", lineNo, line, gap, marker, e._type, lineNo, start, msg)
}
