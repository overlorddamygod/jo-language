package lexer

import "fmt"

type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	ERROR   TokenType = "ERROR"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENTIFIER  = "IDENTIFIER" // add, foobar, x, y, ...
	KEYWORD     = "KEYWORD"    // return, if, else, ...
	INT         = "INT"        // 1343456
	FLOAT       = "FLOAT"      // 12.34
	BOOLEAN     = "BOOLEAN"    // true, false
	STRING      = "STRING"     // "LOL 12312213"
	OPERATOR    = "OPERATOR"
	PUNCTUATION = "PUNCTUATION"
)

const (
	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	BANG        = "!"
	ASTERISK    = "*"
	SLASH       = "/"
	EQ          = "=="
	NOT_EQ      = "!="
	PIPE        = "|"
	AMPERSAND   = "&"
	OR          = "||"
	AND         = "&&"
	LT_EQ       = "<="
	GT_EQ       = ">="
	PERCENT     = "%"
	UNARY_PLUS  = "++"
	UNARY_MINUS = "--"

	PLUS_ASSIGN      = "+="
	MINUS_ASSIGN     = "-="
	BANG_ASSIGN      = "!="
	ASTERISK_ASSIGN  = "*="
	SLASH_ASSIGN     = "/="
	PIPE_ASSIGN      = "|="
	AMPERSAND_ASSIGN = "&="
	AND_ASSIGN       = "&&="
	OR_ASSIGN        = "||="

	// Delimiters
	FULL_STOP = "."
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	DOUBLE_QUOTE = "\""
	SINGLE_QUOTE = "'"

	LT       = "<"
	GT       = ">"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	TRUE  = "true"
	FALSE = "false"
)

type Token struct {
	Type    TokenType
	Literal string
	line    int
	start   int
	end     int
}

func NewToken(tokenType TokenType, literal string) *Token {
	return &Token{Type: tokenType, Literal: literal, line: 0, start: 0, end: 0}
}

func (t *Token) GetLine() int {
	return t.line
}

func (t *Token) Line(line int) *Token {
	t.line = line
	return t
}

func (t *Token) Start(start int) *Token {
	t.start = start
	return t
}
func (t *Token) End(end int) *Token {
	t.end = end
	return t
}
func (t *Token) SetType(_type TokenType) *Token {
	t.Type = _type
	return t
}

func (t *Token) Print() {
	fmt.Printf("%s %s\n", string(t.Type), t.Literal)
}

func IsKeyword(identifier string) bool {
	switch identifier {
	case "return", "if", "let", "else", "struct", "fn", "var", "for", "while", "break", "continue", "true", "false", "null", "int", "string", "char":
		return true
	}
	return false
}

func IsAssignmentOperator(op string) bool {
	switch op {
	case ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN, ASTERISK_ASSIGN, SLASH_ASSIGN, BANG_ASSIGN, PIPE_ASSIGN, AND_ASSIGN, OR_ASSIGN, AMPERSAND_ASSIGN:
		return true
	default:
		return false
	}
}
