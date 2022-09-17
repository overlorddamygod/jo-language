package eval

import (
	"fmt"
	"strconv"

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
	name     string
	_type    string
	Value    string
	FloatVal float64
	IntVal   int64
}

func NewLiteralData(Type, value string) *LiteralData {
	litVal := LiteralData{
		name:  JoLiteral,
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
	return l.Type() == JoInt || l.Type() == JoFloat
}

func (l *LiteralData) IsString() bool {
	return l.Type() == JoString
}

func (l *LiteralData) IsBoolean() bool {
	return l.Type() == JoBoolean
}

func (l *LiteralData) IsNull() bool {
	return l.Type() == JoNull
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

func (l LiteralData) GetBoolean() bool {
	if l.IsNumber() {
		if l.Type() == JoInt {
			return l.IntVal > 0
		}
		return l.FloatVal > 0
	}
	if l.IsString() {
		return true
	}
	if l.IsNull() {
		return false
	}
	return l.Value == "true"
}

func (l LiteralData) GetString() string {
	return l.Value
}

func (l *LiteralData) NodeName() string {
	return l.name
}
func (l *LiteralData) Print() {
	fmt.Println(*l)
}

func BooleanLiteral(boolean bool) LiteralData {
	return *NewLiteralData(JoBoolean, fmt.Sprintf("%v", boolean))
}

func NumberLiteralFloat(val float64) LiteralData {
	return *NewLiteralData(JoFloat, fmt.Sprintf("%f", val))
}

func NumberLiteralInt(val int64) LiteralData {
	return *NewLiteralData(JoInt, fmt.Sprintf("%d", val))
}

func StringLiteral(val string) LiteralData {
	return *NewLiteralData(JoString, val)
}

func NullLiteral() LiteralData {
	return *NewLiteralData(JoNull, JoNull)
}

var parserToLiteralData = map[string]string{
	"NULL":    JoNull,
	"INT":     JoInt,
	"FLOAT":   JoFloat,
	"BOOLEAN": JoBoolean,
	"STRING":  JoString,
}

func LiteralDataFromParserLiteral(li node.LiteralValue) LiteralData {
	JoType, ok := parserToLiteralData[li.Type]

	if !ok {
		panic("UNKNOWN DATA TYPE" + li.Type)
	}
	return *NewLiteralData(JoType, li.Value)
}

func (l *LiteralData) Call(env *Evaluator, name string, arguments []node.Node) (EnvironmentData, error) {
	switch name {
	case "len":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		if l.Type() != JoArray && l.Type() != JoString {
			break
		}
		return NumberLiteralInt(int64(len(l.Value))), nil
	case "type":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		return StringLiteral(l.Type()), nil
	case "getInt":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		if l.Type() != JoInt && l.Type() != JoFloat && l.Type() != JoString {
			break
		}

		intVal, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			return nil, ErrParseInt
		}
		return NumberLiteralInt(int64(intVal)), nil
	case "getFloat":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		if l.Type() != JoInt && l.Type() != JoFloat && l.Type() != JoString {
			break
		}

		floatVal, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			return nil, ErrParseFloat
		}
		return NumberLiteralFloat(floatVal), nil
	case "getString":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		return StringLiteral(l.Value), nil
	}
	return nil, ErrNoMethod(name, l.Type())
}
