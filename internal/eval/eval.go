package eval

import (
	"fmt"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

type Evaluator struct {
	lexer       *L.Lexer
	node        []Node.Node
	global      *Environment
	environment *Environment
	current     Node.Node
}

func NewEvaluator(lexer *L.Lexer, node []Node.Node) *Evaluator {
	env := NewEnvironment()
	return &Evaluator{lexer: lexer, node: node, global: env, environment: env}
}

func Init(src string) {
	_lexer := L.NewLexer(src)

	_, token, err := _lexer.Lex()

	if err != nil {
		//stdio.Io.Println(tokens)
		stdio.Io.Error("[Lexer]\n\n" + JoError.New(_lexer, token, JoError.LexicalError, err.Error()).Error())
		return
	}

	_parser := parser.NewParser(_lexer)

	node, err := _parser.Parse()

	if err != nil {
		stdio.Io.Error("[Parser]\n\n" + err.Error())
		return
	}

	// for _, s := range node {
	// 	s.Print()
	// }

	evaluator := NewEvaluator(_lexer, node)

	_, err = evaluator.Eval()

	if err != nil {
		stdio.Io.Error("[Evaluator]\n\n" + err.Error())
		return
	}
}

func (e *Evaluator) SetLexerNode(lexer *L.Lexer, node []Node.Node) {
	e.lexer = lexer
	e.node = node
}

func NewEvaluatorWithParent(e *Evaluator, parent *Environment) *Evaluator {
	env := NewEnvironmentWithParent(parent)
	return &Evaluator{lexer: e.lexer, node: e.node, global: env, environment: env}
}

func (e *Evaluator) Eval() (EnvironmentData, error) {
	return e.EvalStatements(e.node)
}

func (e *Evaluator) EvalStatements(statements []Node.Node) (EnvironmentData, error) {
	for _, s := range statements {
		data, err := e.EvalStatement(s)
		// fmt.Println("EVALSTATEMENT", s, data, err, err != nil)

		if err != nil {
			return nil, err
		}

		if s.NodeName() == "ReturnStatement" {
			// fmt.Println()
			return data, nil
		}

		if s.NodeName() == "IF" || s.NodeName() == "WHILE" || s.NodeName() == "FOR" || s.NodeName() == "ReturnStatement" {
			if data != nil {
				return data, nil
			}
			// return data, nil
			// fmt.Println("IFFFF", s, data, nil)
			// return data, nil
		}
	}
	return nil, nil
}

func (e *Evaluator) EvalStatement(node Node.Node) (EnvironmentData, error) {
	// e.Eval()
	// fmt.Println("___")
	// fmt.Println("NODE", node.NodeName())
	// node.Print()
	// fmt.Println("___")
	// return nil, nil
	switch node.NodeName() {
	case "VarDecl":
		return e.varDecl(node)
	case "StructDecl":
		return e.structDecl(node)
	case "FunctionDecl":
		return e.functionDecl(node)
	case "ASSIGNMENT":
		return e.assignment(node)
	case "FunctionCall":
		return e.functionCall(node)
	case "IF":
		return e.IfElse(node)
	case "FOR":
		return e.For(node)
	case "WHILE":
		return e.While(node)
	case "Identifier", "BinaryExpression", "GetExpr":
		return e.EvalExpression(node)
	case "ReturnStatement":
		return e.Return(node)
	case "BreakStatement":
		if e.current != nil && (e.current.NodeName() == "FOR" || e.current.NodeName() == "WHILE") {
			return nil, ErrBreak
		}
		return nil, nil
	case "ContinueStatement":
		if e.current != nil && (e.current.NodeName() == "FOR" || e.current.NodeName() == "WHILE") {
			return nil, ErrContinue
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown statement %s", node.NodeName())
	}
}

func (e *Evaluator) begin() {
	e.environment = NewEnvironmentWithParent(e.environment)
}

func (e *Evaluator) end() {
	if e.environment.parent == nil {
		e.environment = e.global
		return
	}
	e.environment = e.environment.parent
}

func (e *Evaluator) NewTokenFromLine(line int) *L.Token {
	return L.NewToken(L.IDENTIFIER, "", line, 0, 0)
}

func (e *Evaluator) NewError(token *L.Token, _type JoError.JoErrorType, message string) error {
	return JoError.New(e.lexer, token, _type, message)
}
