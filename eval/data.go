package eval

import (
	"errors"
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/jo/lexer"
	"github.com/overlorddamygod/jo/parser"
)

type LangData string

var (
	Literal  LangData = "LiteralData"
	Function          = "CallableFunction"
)

type LiteralData struct {
	name           LangData
	_type          string
	Value          string
	NumericalValue float64
}

func NewLiteralData(Type, value string) *LiteralData {
	litVal := LiteralData{
		name:  Literal,
		_type: Type,
		Value: value,
	}
	if Type == L.FLOAT || Type == L.INT {
		litVal.NumericalValue, _ = strconv.ParseFloat(litVal.Value, 32)
	}
	return &litVal
}

func (l LiteralData) Type() string {
	return l._type
}

func (l *LiteralData) IsNumber() bool {
	return l.Type() == L.INT || l.Type() == L.FLOAT
}

func (l *LiteralData) IsString() bool {
	return l.Type() == L.STRING
}

func (l *LiteralData) IsBoolean() bool {
	return l.Type() == L.BOOLEAN
}

func (l *LiteralData) GetNumber() float64 {
	if l.IsBoolean() {
		if l.Value == "true" {
			return 1
		}
		return 0
	}

	if l.IsString() {
		NumericalValue, err := strconv.ParseFloat(l.Value, 32)

		if err != nil {
			return 1
		}
		return NumericalValue
	}
	return l.NumericalValue
}

func (l *LiteralData) GetBoolean() bool {
	if l.IsNumber() {
		return l.GetNumber() > 0
	}
	if l.IsString() {
		return true
	}
	return l.Value == "true"
}

func (l LiteralData) GetString() string {
	return l.Value
}

func (l *LiteralData) NodeName() LangData {
	return l.name
}
func (l *LiteralData) Print() {
	fmt.Println(*l)
}

func BooleanLiteral(boolean bool) LiteralData {
	return *NewLiteralData(L.BOOLEAN, fmt.Sprintf("%v", boolean))
}
func NumberLiteral(val float64) LiteralData {
	return *NewLiteralData(L.FLOAT, fmt.Sprintf("%f", val))
}

func StringLiteral(val string) LiteralData {
	return *NewLiteralData(L.STRING, val)
}

func LiteralDataFromParserLiteral(li parser.LiteralValue) LiteralData {
	return *NewLiteralData(li.Type, li.Value)
}

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
	eval := NewEvaluatorWithParent(e, f.Closure)

	for i, param := range f.FunctionDecl.Params {
		paramId := param.(*parser.Identifier)

		exp, err := e.EvalExpression(arguments[i])

		if err != nil {
			return nil, err
		}
		eval.environment.Define(paramId.Value, exp)
	}

	bodyNodes := f.FunctionDecl.Body.Nodes
	data, err := eval.EvalStatements(bodyNodes)

	if err != nil {
		return nil, err
	}

	return data, nil
}

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
