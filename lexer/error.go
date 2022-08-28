package lexer

import "fmt"

type JoError struct {
	token    *Token
	errorMsg string
}

func (m *JoError) Error() string {
	return m.errorMsg
}

func NewJoError(l *Lexer, token *Token, msg string) *JoError {
	line := l.GetLine(token.line)

	return &JoError{
		token:    token,
		errorMsg: MarkError(line, token.line, token.start, token.end, msg),
	}
}

func MarkError(line string, lineNo int, start int, end int, msg string) string {
	strlen := len(line)
	formatStr := "%s\n "

	for i := 0; i <= strlen; i++ {
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
