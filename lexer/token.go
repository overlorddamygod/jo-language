package lexer

import "fmt"

type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENTIFIER  = "IDENTIFIER" // add, foobar, x, y, ...
	KEYWORD     = "KEYWORD"    // return, if, else, ...
	INT         = "INT"        // 1343456
	FLOAT       = "FLOAT"      // 12.34
	STRING      = "STRING"     // "LOL 12312213"
	OPERATOR    = "OPERATOR"
	PUNCTUATION = "PUNCTUATION"
)

const (
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="
	PERCENT  = "%"

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
)

type Token struct {
	Type    TokenType
	Literal string
	line    int
	col     int
}

func NewToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, line: 0, col: 0}
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
