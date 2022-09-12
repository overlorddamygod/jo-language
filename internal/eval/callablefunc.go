package eval

import (
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

type CallableFunc struct {
	name    string
	_type   string
	Closure *Environment
	Arity   int
	fun     FuncType
}

type FuncType func(*Evaluator, string, []Node.Node) (EnvironmentData, error)

func NewCallableFunc(name string, env *Environment, arity int, fun FuncType) *CallableFunc {
	return &CallableFunc{
		name:    name,
		_type:   Function,
		fun:     fun,
		Arity:   arity,
		Closure: env,
	}
}

func (f CallableFunc) Type() string {
	return f._type
}

func (f CallableFunc) GetString() string {
	return f.name
}

func (f *CallableFunc) Call(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	return f.fun(e, name, arguments)
}
