package eval

import (
	"errors"
	"fmt"

	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

type LangData string

var (
	Literal    LangData = "LiteralData"
	Function            = "CallableFunction"
	StructDecl          = "StructDataDecl"
	Struct              = "StructData"
	JoArray    LangData = "Array"
	Null                = "Null"
)

const INIT_METHOD = "init"

type StructDataDecl struct {
	name       string
	_type      string
	StructDecl Node.StructDeclStatement
	Closure    *Environment
}

func NewStructDataDecl(functionDecl Node.StructDeclStatement, env *Environment) *StructDataDecl {
	return &StructDataDecl{
		name:       StructDecl,
		_type:      StructDecl,
		StructDecl: functionDecl,
		Closure:    env,
	}
}

func (s *StructDataDecl) Initialize(e *Evaluator, args []Node.Node) (*StructData, error) {
	env := NewEnvironmentWithParent(s.Closure)

	methods := s.StructDecl.Methods

	structData := &StructData{
		name:       Struct,
		_type:      Struct,
		StructDecl: s.StructDecl,
		env:        env,
	}
	initFound := false

	// TODO: Declare methods in StructDataDecl only ??
	for _, method := range methods {
		id := method.Identifier.(*Node.Identifier)
		if id.Value == INIT_METHOD {
			if initFound {
				return nil, errors.New("init (constructor) method already defined")
			}
			initFound = true
		}
		env.Define(id.Value, NewCallableFunction(method, env, structData))
	}
	env.Define("self", structData)

	env.Define("type", NewCallableFunc("type", env, 0, func(e *Evaluator, name string, n []Node.Node) (EnvironmentData, error) {
		id := s.StructDecl.Identifier.(*Node.Identifier)
		return StringLiteral(id.Value), nil
	}))

	d, err := structData.Get(INIT_METHOD)

	if err != nil {
		return structData, nil
	}

	de := d.(*CallableFunction)

	if len(args) != len(de.FunctionDecl.Params) {
		return nil, errors.New("failed to initialize struct. arguments length does not match")
	}

	structData.Call(e, INIT_METHOD, args)

	env.Remove(INIT_METHOD)
	// env.Define("sad", NumberLiteral(69))

	return structData, nil
}

func (s StructDataDecl) Type() string {
	return s._type
}

func (s StructDataDecl) GetString() string {
	return fmt.Sprintf("[structDecl %s]", s.StructDecl.Identifier.(*Node.Identifier).Value)
}

type StructData struct {
	name       string
	_type      string
	StructDecl Node.StructDeclStatement
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

	// TODO: Declare methods in StructDataDecl only ??
	for _, method := range methods {
		id := method.Identifier.(*Node.Identifier)
		env.Define(id.Value, NewCallableFunction(method, env, structData))
	}
	env.Define("self", structData)

	// d, err := structData.Get("init")

	// if err != nil {
	// 	return structData
	// }

	// init, ok := d.(*CallableFunction)

	// if ok {
	// 	structData.Call(e, "init", args)
	// }

	// env.Define("sad", NumberLiteral(69))

	return structData
}

func NewNativeStruct(env *Environment, name string, methods []*CallableFunc) *StructData {
	env = NewEnvironmentWithParent(env)

	structData := &StructData{
		name:  name,
		_type: Struct,
		env:   env,
	}

	env.Define("self", structData)

	for _, method := range methods {
		env.Define(method.name, method)
	}

	return structData
}

func (s *StructData) Get(key string) (EnvironmentData, error) {
	return s.env.Get(key)
}

// not used
func (s *StructData) Call(e *Evaluator, funcName string, args []Node.Node) (EnvironmentData, error) {
	data, err := s.Get(funcName)

	if err != nil {
		return nil, err
	}

	fun, ok := data.(*CallableFunction)

	if ok {
		// return nil, L.NewJoError(e.lexer, nil, "not a function")
		return fun.Call(e, funcName, args)
	}

	f, ok := data.(*CallableFunc)

	if ok {
		// return nil, L.NewJoError(e.lexer, nil, "not a function")
		return f.Call(e, funcName, args)
	}
	return nil, errors.New("not a function")
}

func (s *StructData) CallWithEnvData(e *Evaluator, funcName string, args []EnvironmentData) (EnvironmentData, error) {
	data, err := s.Get(funcName)

	if err != nil {
		return nil, err
	}

	fun, ok := data.(*CallableFunction)
	if ok {
		return fun.CallWithEnvData(e, funcName, args)
	}

	// f, ok := data.(*CallableFunc)

	// if ok {
	// 	// return nil, L.NewJoError(e.lexer, nil, "not a function")
	// 	return f.Call(e, funcName, args)
	// }
	return nil, errors.New("not a function")
}

func (f StructData) Type() string {
	return f._type
}

func (f *StructData) SetName(name string) {
	f.name = name
}

func (s StructData) GetString() string {
	if s.name != Struct {
		return fmt.Sprintf("[struct %s]", s.name)
	}
	return fmt.Sprintf("[struct %s]", s.StructDecl.Identifier.(*Node.Identifier).Value)
}
