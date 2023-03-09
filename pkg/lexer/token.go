package lexer

import "fmt"

type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	ERROR   TokenType = "ERROR"
	EOF     TokenType = "EOF"

	WHITESPACE = "WHITESPACE"
	NEWLINE    = "NEWLINE"

	// Identifiers + literals
	IDENTIFIER  = "IDENTIFIER" // add, foobar, x, y, ...
	KEYWORD     = "KEYWORD"    // return, if, else, ...
	INT         = "INT"        // 1343456
	FLOAT       = "FLOAT"      // 12.34
	BOOLEAN     = "BOOLEAN"    // true, false
	NULL        = "NULL"
	STRING      = "STRING" // "LOL 12312213"
	OPERATOR    = "OPERATOR"
	PUNCTUATION = "PUNCTUATION"
	COMMENT     = "COMMENT"
)

var RegexTokenMap = [][]string{
	{"^\\n", NEWLINE},
	{"^(//.*)", COMMENT},
	{"^\\s+", WHITESPACE},
	{"^\\d+\\.\\d+", FLOAT},
	{"^\\d+", INT},
	{"^[a-zA-Z_0-9]+", IDENTIFIER},
	{"^(\"([^\\\"]|\\.)*\")", STRING},
	{"^\\.", OPERATOR},
	{"^==", OPERATOR}, // Equality
	{"^[\\(\\)\\{\\}\\[\\];,;:]", PUNCTUATION},
	{"^(&&|\\|\\||[+\\-*\\/%!\\|&])?=", OPERATOR}, // Assignment
	{"^((\\+\\+|\\-\\-)|[+\\-*\\/%])", OPERATOR},  // Arithmetic
	{"^(&&|\\|\\||!)", OPERATOR},                  // Logical
	{"^([<>!=]=)|([><])", OPERATOR},               // Relational
	{"^\\||&", OPERATOR},                          // Bitwise
}

const (
	// Arithmetic Operators
	PLUS        = "+"
	MINUS       = "-"
	ASTERISK    = "*"
	SLASH       = "/"
	PERCENT     = "%"
	UNARY_PLUS  = "+"
	UNARY_MINUS = "-"

	// Relational Operators
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LT_EQ  = "<="
	GT_EQ  = ">="

	// Logical Operators
	AND  = "&&"
	OR   = "||"
	BANG = "!"

	// Bitwise Operators
	PIPE      = "|"
	AMPERSAND = "&"

	// ASSIGNMENT Operators
	ASSIGN           = "="
	PLUS_ASSIGN      = "+="
	MINUS_ASSIGN     = "-="
	ASTERISK_ASSIGN  = "*="
	SLASH_ASSIGN     = "/="
	PERCENT_ASSIGN   = "%="
	BANG_ASSIGN      = "!="
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

func NewTokenWithoutLine(tokenType TokenType, literal string) *Token {
	return &Token{Type: tokenType, Literal: literal, line: 0, start: 0, end: 0}
}

func NewToken(tokenType TokenType, literal string, line, start, end int) *Token {
	return &Token{
		Type:    tokenType,
		Literal: literal,
		line:    line,
		start:   start,
		end:     end,
	}
}

func (t *Token) GetLine() int {
	return t.line
}

func (t *Token) GetEnd() int {
	return t.end
}

func (t *Token) GetStart() int {
	return t.start
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
	fmt.Printf("%s %s %d %d %d\n", string(t.Type), t.Literal, t.line, t.start, t.end)
}

func IsKeyword(identifier string) bool {
	switch identifier {
	case "import", "export", "return", "if", "let", "else", "struct", "fn", "var", "for", "while", "break", "continue", "true", "false", "null", "string", "char", "switch", "case", "default", "try", "catch", "throw":
		return true
	}
	return false
}

func IsAssignmentOperator(op string) bool {
	switch op {
	case ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN, ASTERISK_ASSIGN, SLASH_ASSIGN, BANG_ASSIGN, PIPE_ASSIGN, AND_ASSIGN, OR_ASSIGN, AMPERSAND_ASSIGN, PERCENT_ASSIGN:
		return true
	default:
		return false
	}
}
