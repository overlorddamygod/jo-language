package eval

import (
	"errors"
	"fmt"

	"github.com/overlorddamygod/jo/pkg/parser"
)

type LangData string

var (
	Literal    LangData = "LiteralData"
	Function            = "CallableFunction"
	StructDecl          = "StructDataDecl"
	Struct              = "StructData"
)

type StructDataDecl struct {
	name       string
	_type      string
	StructDecl parser.StructDeclStatement
	Closure    *Environment
}

func NewStructDataDecl(functionDecl parser.StructDeclStatement, env *Environment) *StructDataDecl {
	return &StructDataDecl{
		name:       StructDecl,
		_type:      StructDecl,
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

func NewStructData(structDecl StructDataDecl) *StructData {
	env := NewEnvironmentWithParent(structDecl.Closure)

	methods := structDecl.StructDecl.Methods

	structData := &StructData{
		name:       Struct,
		_type:      Struct,
		StructDecl: structDecl.StructDecl,
		env:        env,
	}

	for _, method := range methods {
		id := method.Identifier.(*parser.Identifier)
		env.Define(id.Value, NewCallableFunction(method, env, structData))
	}

	env.Define("self", structData)

	// env.Define("sad", NumberLiteral(69))

	return structData
}

func (s *StructData) Get(key string) (EnvironmentData, error) {
	return s.env.Get(key)
}

// not used
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
