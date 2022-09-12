package node

const (
	IF                = "IF"
	FOR               = "FOR"
	WHILE             = "WHILE"
	BREAK             = "BREAK"
	BLOCK             = "BLOCK"
	CONTINUE          = "CONTINUE"
	STRUCT_DECL       = "STRUCT_DECL"
	VAR_DECL          = "VAR_DECL"
	FUNCTION_CALL     = "FUNCTION_CALL"
	FUNCTION_DECL     = "FUNCTION_DECL"
	RETURN            = "RETURN"
	BINARY_EXPRESSION = "BINARY_EXPRESSION"
	UNARY_EXPRESSION  = "UNARY_EXPRESSION"
	LITERAL_VALUE     = "LITERAL_VALUE"
	IDENTIFIER        = "IDENTIFIER"
	GET_EXPR          = "GET_EXPR"
	ASSIGNMENT        = "ASSIGNMENT"
	CONDITION_BLOCK   = "CONDITION_BLOCK"
	ARRAY             = "ARRAY"
)

type Node interface {
	NodeName() string
	Print()
	GetLine() int
}
