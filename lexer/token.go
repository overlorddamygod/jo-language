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
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	EQ        = "=="
	NOT_EQ    = "!="
	PIPE      = "|"
	AMPERSAND = "&"
	OR        = "||"
	AND       = "&&"
	LT_EQ     = "<="
	GT_EQ     = ">="
	PERCENT   = "%"

	// Delimiters
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
	case "return", "if", "else", "func", "var", "for", "while", "break", "continue", "true", "false", "nil", "int", "string", "char", "class", "const":
		return true
	}
	return false
}
