package eval

import (
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/jo/lexer"
	"github.com/overlorddamygod/jo/parser"
)

type LiteralData struct {
	name           string
	_type          string
	Value          string
	NumericalValue float64
}

func NewLiteralData(Type, value string) *LiteralData {
	litVal := LiteralData{
		name:  "LiteralData",
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

func (l *LiteralData) NodeName() string {
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

type CallableFunction struct {
	name         string
	_type        string
	FunctionDecl parser.FunctionDeclStatement
	env          *Environment
	parent       *Environment
}

func NewCallableFunction(functionDecl parser.FunctionDeclStatement) *CallableFunction {
	return &CallableFunction{
		name:         "CallableFunction",
		_type:        "CallableFunction",
		FunctionDecl: functionDecl,
	}
}

func (f CallableFunction) Type() string {
	return f._type
}

func (f *CallableFunction) Exec(e *Evaluator, arguments []parser.Node) (EnvironmentData, error) {
	// e.environment.Print()
	// f.parent = e.environment
	// eval := NewEvaluatorWithParent(e, e.global)
	// e.begin()
	// e.environment = NewEnvironmentWithParent(e.global)
	// f.env = e.environment

	// env := e.environment
	// fmt.Println("FUNC START")
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
	eval := NewEvaluatorWithParent(e, e.global)
	// e.global.Print()
	// e.environment.parent.Print()

	for i, param := range f.FunctionDecl.Params {
		paramId := param.(*parser.Identifier)

		exp, err := e.EvalExpression(arguments[i])

		if err != nil {
			return nil, err
		}
		eval.environment.Define(paramId.Value, exp)
	}

	// e.environment.Print()
	bodyNodes := f.FunctionDecl.Body.Nodes
	data, err := eval.EvalStatements(bodyNodes)

	// fmt.Println("EEEE", data, err)
	if err != nil {
		return nil, err
	}
	// fmt.Println("PREVVVVVV")
	// e.environment.Print()

	// eval.end()

	// e.environment = f.parent

	// fmt.Println("NEXTTTT")

	// e.environment.Print()
	// fmt.Println("FUNC END")

	return data, nil
}
