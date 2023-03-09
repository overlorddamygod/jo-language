package joerror

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

type JoRuntimeError struct {
	Type     JoErrorType
	Token    *lexer.Token
	ErrorMsg error
	Parent   *error
}

type JoErrorType string

var (
	DefaultError   JoErrorType = "Error"
	LexicalError   JoErrorType = "LexicalError"
	SyntaxError    JoErrorType = "SyntaxError"
	ReferenceError JoErrorType = "ReferenceError"
)

func (m *JoRuntimeError) Error() string {
	errStr := fmt.Sprintf("  File \"filename\", line: %d, in %s\n", m.Token.GetLine(), m.Token.Literal)
	// fmt.Println("ZZZ", m.Token)
	if m.Parent != nil {
		errStr += (*m.Parent).Error()
	}
	if m.ErrorMsg != nil {
		errStr += fmt.Sprintf("%s: %s", string(m.Type), m.ErrorMsg.Error())
	}

	return errStr
	// return m.errorMsg.Error()
}

func NewRuntimeError(l *lexer.Lexer, token *lexer.Token, _type JoErrorType, err interface{}) *JoRuntimeError {
	joErr := &JoRuntimeError{
		Token: token,
		Type:  _type,
		// parent: err,
		// errorMsg: err,
	}
	// fmt.Println("NNNNN:", err, token, _type)
	errStr, ok := err.(string)

	if ok {
		joErr.ErrorMsg = errors.New(errStr)
	} else {
		errErr := err.(error)
		joErr.Parent = &errErr
	}

	return joErr
}

// func (e *JoError)

type JoError struct {
	Type     JoErrorType
	Token    *lexer.Token
	ErrorMsg error
}

func (m *JoError) Error() string {
	return m.ErrorMsg.Error()
}

func New(l *lexer.Lexer, token *lexer.Token, _type JoErrorType, msg string) *JoError {
	var err error = fmt.Errorf("[%s] %s", _type, msg)

	joErr := &JoError{
		Token:    token,
		Type:     _type,
		ErrorMsg: err,
	}
	if token != nil {
		line, err := l.GetLine(token.GetLine())
		if err != nil {
			fmt.Println(err.Error())
		}
		err = errors.New(joErr.MarkError(line, token.GetLine(), token.GetStart(), token.GetEnd(), msg))
		joErr.ErrorMsg = err
	}
	return joErr
}

func (e *JoError) MarkError(line string, lineNo int, start int, end int, msg string) string {
	lineNoLength := len(strconv.Itoa(lineNo))

	gap := strings.Repeat(" ", 2+lineNoLength+start)
	marker := strings.Repeat("^", end-start+1)

	return fmt.Sprintf("%d | %s   \n%s%s   \n-- [ %s ] Line: %d Col: %d : %s\n", lineNo, line, gap, marker, e.Type, lineNo, start, msg)
}
