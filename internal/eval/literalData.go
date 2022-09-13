package eval

import (
	"errors"
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

// type Number struct {
// 	FloatVal float64
// 	IntVal   int64
// }

// func NewNumber() Number{
// 	return
// }

type LiteralData struct {
	name     LangData
	_type    string
	Value    string
	FloatVal float64
	IntVal   int64
}

func NewLiteralData(Type, value string) *LiteralData {
	litVal := LiteralData{
		name:  Literal,
		_type: Type,
		Value: value,
	}
	litVal.FloatVal, _ = strconv.ParseFloat(litVal.Value, 64)
	litVal.IntVal, _ = strconv.ParseInt(litVal.Value, 10, 64)
	// fmt.Println(litVal)
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

func (l *LiteralData) GetNumber() (int64, float64) {
	if l.IsBoolean() {
		if l.Value == "true" {
			return 1, 1
		}
		return 0, 0
	}

	if l.IsString() {
		floatVal, err := strconv.ParseFloat(l.Value, 64)

		if err != nil {
			return 1, 1
		}
		intVal, err := strconv.ParseInt(l.Value, 10, 64)

		if err != nil {
			return 1, 1
		}
		return intVal, floatVal
	}
	return l.IntVal, l.FloatVal
}

func (l *LiteralData) GetBoolean() bool {
	if l.IsNumber() {
		if l.Type() == L.INT {
			return l.IntVal > 0
		}
		return l.FloatVal > 0
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

func NumberLiteralFloat(val float64) LiteralData {
	return *NewLiteralData(L.FLOAT, fmt.Sprintf("%f", val))
}

func NumberLiteralInt(val int64) LiteralData {
	return *NewLiteralData(L.INT, fmt.Sprintf("%d", val))
}

func StringLiteral(val string) LiteralData {
	return *NewLiteralData(L.STRING, val)
}

func LiteralDataFromParserLiteral(li node.LiteralValue) LiteralData {
	return *NewLiteralData(li.Type, li.Value)
}

func (l *LiteralData) Call(env *Evaluator, name string, arguments []node.Node) (EnvironmentData, error) {
	switch name {
	case "len":
		return NumberLiteralInt(int64(len(l.Value))), nil
	case "type":
		if len(arguments) != 0 {
			return nil, errors.New("argument length must be 0")
		}
		return StringLiteral(l._type), nil
	case "getInt":
		if len(arguments) != 0 {
			return nil, errors.New("argument length must be 0")
		}

		intVal, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			return nil, errors.New("cannot parse to int")
		}
		return NumberLiteralInt(int64(intVal)), nil
	case "getFloat":
		if len(arguments) != 0 {
			return nil, errors.New("argument length must be 0")
		}

		floatVal, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			return nil, errors.New("cannot parse to float")
		}
		return NumberLiteralFloat(floatVal), nil
	case "getString":
		if len(arguments) != 0 {
			return nil, errors.New("argument length must be 0")
		}

		return StringLiteral(l.Value), nil
	}
	return nil, errors.New("no method")
}
