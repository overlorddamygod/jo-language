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
	evaluator.LoadNative()

	_, err = evaluator.Eval()

	if err != nil {
		stdio.Io.Error("[Evaluator]\n\n" + err.Error())
		return
	}
}

func (e *Evaluator) LoadNative() {
	e.environment.Define("print", NewCallableFunc("print", e.global, -1, Print))
	e.environment.Define("input", NewCallableFunc("input", e.global, 1, Input))
	e.environment.Define("math", Math(e))
}

func (e *Evaluator) SetLexerNode(lexer *L.Lexer, node []Node.Node) {
	e.lexer = lexer
	e.node = node
}
func (e *Evaluator) Env() *Environment {
	return e.environment
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

		if s.NodeName() == Node.IF || s.NodeName() == Node.WHILE || s.NodeName() == Node.FOR || s.NodeName() == Node.RETURN || s.NodeName() == Node.SWITCH {
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
	case Node.VAR_DECL:
		return e.varDecl(node)
	case Node.STRUCT_DECL:
		return e.structDecl(node)
	case Node.FUNCTION_DECL:
		return e.functionDecl(node)
	case Node.ASSIGNMENT:
		return e.assignment(node)
	case Node.FUNCTION_CALL:
		return e.functionCall(node)
	case Node.SWITCH:
		return e.Switch(node)
	case Node.IF:
		return e.IfElse(node)
	case Node.FOR:
		return e.For(node)
	case Node.WHILE:
		return e.While(node)
	case Node.BLOCK:
		blockStmt := node.(*Node.Block)
		e.begin()

		_, err := e.EvalStatements(blockStmt.Nodes)

		e.end()

		if err != nil {
			return nil, err
		}
		return nil, nil
	case Node.IDENTIFIER, Node.BINARY_EXPRESSION, Node.GET_EXPR:
		return e.EvalExpression(node)
	case Node.RETURN:
		return e.Return(node)
	case Node.BREAK:
		if e.current != nil && (e.current.NodeName() == Node.FOR || e.current.NodeName() == Node.WHILE || e.current.NodeName() == Node.SWITCH) {
			return nil, ErrBreak
		}
		return nil, nil
	case Node.CONTINUE:
		if e.current != nil && (e.current.NodeName() == Node.FOR || e.current.NodeName() == Node.WHILE) {
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
