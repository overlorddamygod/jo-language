package eval

import (
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
	if Type == L.INT || Type == L.FLOAT {
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
		return 1
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

func (l *LiteralData) GetString() string {
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
