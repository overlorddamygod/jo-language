package eval

import (
	"errors"
	"fmt"

	"github.com/overlorddamygod/jo/pkg/parser"
)

type LangData string

var (
	Literal  LangData = "LiteralData"
	Function          = "CallableFunction"
)

type StructDataDecl struct {
	name       string
	_type      string
	StructDecl parser.StructDeclStatement
	Closure    *Environment
}

func NewStructDataDecl(functionDecl parser.StructDeclStatement, env *Environment) *StructDataDecl {
	return &StructDataDecl{
		name:       "StructDataDecl",
		_type:      "StructDataDecl",
		StructDecl: functionDecl,
		Closure:    env,
	}
}

func (s StructDataDecl) Type() string {
	return s._type
}

func (s StructDataDecl) GetString() string {
	return fmt.Sprintf("[structDecl %s]", s.StructDecl.Identifier.(*parser.Identifier).Value)
}

type StructData struct {
	name       string
	_type      string
	StructDecl parser.StructDeclStatement
	env        *Environment
}

func NewStructData(structDecl StructDataDecl, env *Environment) *StructData {
	env = NewEnvironmentWithParent(env)

	methods := structDecl.StructDecl.Methods

	for _, method := range methods {
		id := method.Identifier.(*parser.Identifier)
		env.Define(id.Value, NewCallableFunction(method, env))
	}

	return &StructData{
		name:       "StructData",
		_type:      "StructData",
		StructDecl: structDecl.StructDecl,
		env:        env,
	}
}

func (s *StructData) Get(key string) (EnvironmentData, error) {
	return s.env.Get(key)
}

func (s *StructData) Call(funcName string, e *Evaluator, args []parser.Node) (EnvironmentData, error) {

	data, err := s.Get(funcName)

	if err != nil {
		return nil, err
	}

	fun, ok := data.(*CallableFunction)

	if !ok {
		// return nil, L.NewJoError(e.lexer, nil, "not a function")
		return nil, errors.New("not a function")
	}

	return fun.Call(e, args)
}

func (f StructData) Type() string {
	return f._type
}

func (s StructData) GetString() string {
	return fmt.Sprintf("[struct %s]", s.StructDecl.Identifier.(*parser.Identifier).Value)
}
