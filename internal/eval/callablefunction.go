package eval

import (
	"errors"
	"fmt"

	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

type Callable interface {
	Call(e *Evaluator, arguments []Node.Node) (EnvironmentData, error)
}

type CallableFunction struct {
	name         string
	_type        string
	FunctionDecl Node.FunctionDeclStatement
	Closure      *Environment
	parent       *StructData
}

func NewCallableFunction(functionDecl Node.FunctionDeclStatement, env *Environment, parent *StructData) *CallableFunction {
	return &CallableFunction{
		name:         Function,
		_type:        Function,
		FunctionDecl: functionDecl,
		parent:       parent,
		Closure:      env,
	}
}

func (f CallableFunction) Type() string {
	return f._type
}

func (f *CallableFunction) GetString() string {
	fName := f.FunctionDecl.Identifier.(*Node.Identifier).Value
	if f.parent == nil {
		return fmt.Sprintf("[function %s]", fName)
	}
	structName := f.parent.StructDecl.Identifier.(*Node.Identifier).Value
	return fmt.Sprintf("[method %s.%s]", structName, fName)
}

func (f *CallableFunction) Call(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	paramsLen := len(f.FunctionDecl.Params)
	argsLen := len(arguments)

	if argsLen > paramsLen {
		// iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, errors.New("Arg length greater than params length")
	}

	if argsLen < paramsLen {
		// iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, errors.New("Arg length less than params length")
	}
	evaluator := NewEvaluatorWithParent(e, f.Closure)
	evaluator.FunctionScope = e.FunctionScope

	for i, param := range f.FunctionDecl.Params {
		paramId := param.(*Node.Identifier)

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
	if data != nil {
		return data, nil
	}
	return NullLiteral(), nil
}

func (f *CallableFunction) CallWithEnvData(e *Evaluator, name string, arguments []EnvironmentData) (EnvironmentData, error) {
	paramsLen := len(f.FunctionDecl.Params)
	argsLen := len(arguments)

	if argsLen > paramsLen {
		// iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, errors.New("Arg length greater than params length")
	}

	if argsLen < paramsLen {
		// iden := f.FunctionDecl.Identifier.(*parser.Identifier)
		return nil, errors.New("Arg length less than params length")
	}
	evaluator := NewEvaluatorWithParent(e, f.Closure)
	evaluator.FunctionScope = e.FunctionScope

	for i, param := range f.FunctionDecl.Params {
		paramId := param.(*Node.Identifier)

		evaluator.environment.Define(paramId.Value, arguments[i])
	}

	bodyNodes := f.FunctionDecl.Body.Nodes
	data, err := evaluator.EvalStatements(bodyNodes)
	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}
	return NullLiteral(), nil
}
