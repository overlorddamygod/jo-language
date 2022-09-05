package eval

import (
	"fmt"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser"
)

type Callable interface {
	Call(e *Evaluator, arguments []parser.Node) (EnvironmentData, error)
}

type CallableFunction struct {
	name         string
	_type        string
	FunctionDecl parser.FunctionDeclStatement
	Closure      *Environment
}

func NewCallableFunction(functionDecl parser.FunctionDeclStatement, env *Environment) *CallableFunction {
	return &CallableFunction{
		name:         Function,
		_type:        Function,
		FunctionDecl: functionDecl,
		Closure:      env,
	}
}

func (f CallableFunction) Type() string {
	return f._type
}

func (f *CallableFunction) GetString() string {
	return fmt.Sprintf("[function %s]", f.FunctionDecl.Identifier.(*parser.Identifier).Value)
}

func (f *CallableFunction) Call(e *Evaluator, arguments []parser.Node) (EnvironmentData, error) {
	paramsLen := len(f.FunctionDecl.Params)
	argsLen := len(arguments)

	if argsLen > paramsLen {
		iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, L.NewJoError(e.lexer, iden.Token, "Arg length greater than params length")
	}

	if argsLen < paramsLen {
		iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, L.NewJoError(e.lexer, iden.Token, "Arg length less than params length")
	}
	evaluator := NewEvaluatorWithParent(e, f.Closure)

	for i, param := range f.FunctionDecl.Params {
		paramId := param.(*parser.Identifier)

		exp, err := e.EvalExpression(arguments[i])

		if err != nil {
			return nil, err
		}
		evaluator.environment.Define(paramId.Value, exp)
	}

	bodyNodes := f.FunctionDecl.Body.Nodes
	data, err := evaluator.EvalStatements(bodyNodes)

	if err != nil {
		return nil, err
	}

	return data, nil
}
