package eval

import (
	"errors"
	"math"
	"math/rand"
	"time"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

var ErrArgNumber = errors.New("arg must be of type number")

func Math(e *Evaluator) *StructData {
	var methods []*CallableFunc = []*CallableFunc{
		NewCallableFunc("random", e.Env(), 0, random),
		NewCallableFunc("pow", e.Env(), 2, pow),
		NewCallableFunc("exp", e.Env(), 1, exp),

		NewCallableFunc("log", e.Env(), 1, log),
		NewCallableFunc("log2", e.Env(), 1, log2),
		NewCallableFunc("log10", e.Env(), 1, log10),
		NewCallableFunc("sqrt", e.Env(), 1, sqrt),
		NewCallableFunc("abs", e.Env(), 1, abs),

		NewCallableFunc("sin", e.Env(), 1, sin),
		NewCallableFunc("cos", e.Env(), 1, cos),
		NewCallableFunc("tan", e.Env(), 1, tan),

		NewCallableFunc("round", e.Env(), 1, round),
		NewCallableFunc("ceil", e.Env(), 1, ceil),
		NewCallableFunc("floor", e.Env(), 1, floor),
	}

	nativeStruct := NewNativeStruct(e.Env(), "math", methods)

	nativeStruct.env.DefineOne("pi", NumberLiteralFloat(math.Pi))
	nativeStruct.env.DefineOne("e", NumberLiteralFloat(math.E))

	return nativeStruct
}

func random(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	rand.Seed(time.Now().UnixNano())
	float := rand.Float64()
	return NumberLiteralFloat(float), nil
}

func sin(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Sin(val.FloatVal)), nil
}

func cos(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Cos(val.FloatVal)), nil
}

func tan(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Tan(val.FloatVal)), nil
}

func round(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralInt(int64(math.Round(val.FloatVal))), nil
}

func ceil(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralInt(int64(math.Ceil(val.FloatVal))), nil
}

func floor(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralInt(int64(math.Floor(val.FloatVal))), nil
}

func log(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Log(val.FloatVal)), nil
}
func log2(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Log2(val.FloatVal)), nil
}
func log10(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}
	return NumberLiteralFloat(math.Log10(val.FloatVal)), nil
}

func sqrt(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralFloat(math.Sqrt(val.FloatVal)), nil
}
func abs(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralFloat(math.Abs(val.FloatVal)), nil
}
func exp(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	val, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	return NumberLiteralFloat(math.Exp(val.FloatVal)), nil
}
func pow(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	x, err := getNumber(e, arguments[0])

	if err != nil {
		return nil, err
	}

	y, err := getNumber(e, arguments[1])

	if err != nil {
		return nil, err
	}

	return NumberLiteralFloat(math.Pow(x.FloatVal, y.FloatVal)), nil
}

func getNumber(e *Evaluator, node Node.Node) (*LiteralData, error) {
	data, err := e.EvalExpression(node)

	if err != nil {
		return nil, err
	}

	lit, ok := data.(LiteralData)

	if !ok {
		return nil, ErrArgNumber
	}

	if lit.Type() != L.INT && lit.Type() != L.FLOAT {
		return nil, ErrArgNumber
	}

	return &lit, nil
}
