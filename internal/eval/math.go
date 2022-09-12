package eval

import (
	"math/rand"
	"strconv"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func MathRandomInt(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	return NewLiteralData(L.INT, strconv.Itoa(rand.Int())), nil
}

func MathRandomFloat(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	float := rand.Float64()
	return NumberLiteral(float), nil
}

func MathRand(e *Evaluator) *StructData {
	var methods []*CallableFunc = []*CallableFunc{
		NewCallableFunc("int", e.environment, 0, MathRandomInt),
		NewCallableFunc("float", e.environment, 0, MathRandomFloat),
	}

	return NewNativeStruct(e.environment, "rand", methods)
}
