package eval

import (
	"errors"
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser/node"
)

type Array struct {
	name  LangData
	_type string
	data  []EnvironmentData
}

func NewArray(data []EnvironmentData) *Array {
	return &Array{
		name:  JoArray,
		_type: string(JoArray),
		data:  data,
	}
}

func (a Array) Type() string {
	return a._type
}

func (a Array) GetString() string {
	output := "[ "

	length := len(a.data)
	for i, val := range a.data {
		output += val.GetString()
		if i == length-1 {
			output += " "
		} else {
			output += ", "
		}
	}
	output += "]"
	return output
}

func (a *Array) Call(e *Evaluator, name string, arguments []node.Node) (EnvironmentData, error) {
	switch name {
	case "len":
		return NewLiteralData(L.INT, strconv.Itoa(len(a.data))), nil
	case "get":
		if len(arguments) != 1 {
			return nil, errors.New("argument length must be 1")
		}

		index, err := e.EvalExpression(arguments[0])

		if err != nil {
			return nil, err
		}
		if index.Type() != "INT" {
			return nil, errors.New("argument must be of type int")
		}

		i, _ := index.(LiteralData)

		indexInt := int(i.NumericalValue)

		if indexInt < 0 || indexInt > len(a.data)-1 {
			return nil, errors.New("index out of bound")
		}

		data := a.data[indexInt]

		return data, nil
	case "push":
		if len(arguments) == 0 {
			return nil, errors.New("argument length must be greater than 0")
		}

		for _, args := range arguments {
			data, err := e.EvalExpression(args)

			if err != nil {
				return nil, err
			}

			a.data = append(a.data, data)
		}
		return nil, nil
	case "set":
		if len(arguments) != 2 {
			return nil, errors.New("argument length must be 2")
		}
		index, err := e.EvalExpression(arguments[0])

		if err != nil {
			return nil, err
		}

		if index.Type() != "INT" {
			return nil, errors.New("argument must be of type int")
		}

		i, _ := index.(LiteralData)

		indexInt := int(i.NumericalValue)

		if indexInt < 0 || indexInt > len(a.data)-1 {
			return nil, errors.New("index out of bound")
		}

		val, err := e.EvalExpression(arguments[1])

		if err != nil {
			return nil, err
		}

		a.data[indexInt] = val
		return nil, nil
	case "type":
		if len(arguments) != 0 {
			return nil, errors.New("argument length must be 0")
		}
		return StringLiteral("Array"), nil
	}
	return nil, fmt.Errorf("no method `%s`", name)
}
