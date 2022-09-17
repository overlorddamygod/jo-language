package eval

import (
	"errors"

	"github.com/overlorddamygod/jo/pkg/parser/node"
)

type Array struct {
	name  string
	_type string
	data  []EnvironmentData
}

func NewArray(data []EnvironmentData) *Array {
	return &Array{
		name:  JoArray,
		_type: JoArray,
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

func (a Array) GetBoolean() bool {
	return true
}

func (a *Array) Call(e *Evaluator, name string, arguments []node.Node) (EnvironmentData, error) {
	switch name {
	case "len":
		return NumberLiteralInt(int64(len(a.data))), nil
	case "get":
		if _, err := expectArgLength(arguments, 1); err != nil {
			return nil, err
		}

		index, err := getArg(e, JoInt, arguments[0])

		if err != nil {
			return nil, err
		}

		indexLit, _ := index.(LiteralData)
		indexInt := int(indexLit.IntVal)

		if indexInt < 0 || indexInt > len(a.data)-1 {
			return nil, ErrIndexOutofBound
		}

		data := a.data[indexInt]

		return data, nil
	case "push":
		if _, err := expectArgLength(arguments, -1); err != nil {
			return nil, err
		}

		for _, args := range arguments {
			data, err := e.EvalExpression(args)

			if err != nil {
				return nil, err
			}

			a.data = append(a.data, data)
		}
		return NullLiteral(), nil
	case "pop":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}

		if len(a.data) == 0 {
			return nil, errors.New("array is empty")
		}
		last := a.data[len(a.data)-1]
		a.data = a.data[:len(a.data)-1]

		return last, nil
	case "set":
		if _, err := expectArgLength(arguments, 2); err != nil {
			return nil, err
		}

		index, err := getArg(e, JoInt, arguments[0])

		if err != nil {
			return nil, err
		}

		indexLit, _ := index.(LiteralData)
		indexInt := int(indexLit.IntVal)

		if indexInt < 0 || indexInt > len(a.data)-1 {
			return nil, ErrIndexOutofBound
		}

		val, err := e.EvalExpression(arguments[1])

		if err != nil {
			return nil, err
		}

		a.data[indexInt] = val
		return NullLiteral(), nil
	case "contains":
		if _, err := expectArgLength(arguments, 1); err != nil {
			return nil, err
		}

		data, err := e.EvalExpression(arguments[0])

		if err != nil {
			return nil, err
		}

		for _, d := range a.data {
			if d == data {
				return BooleanLiteral(true), nil
			}
		}

		return BooleanLiteral(false), nil
	case "join":
		if _, err := expectArgLength(arguments, 1); err != nil {
			return nil, err
		}
		joinStr, err := getArg(e, JoString, arguments[0])

		if err != nil {
			return nil, err
		}

		joinStrLit, _ := joinStr.(LiteralData)

		str := ""

		for i, val := range a.data {
			str += val.GetString()

			if i != len(a.data)-1 {
				str += joinStrLit.GetString()
			}
		}

		return StringLiteral(str), nil
	case "slice":
		if _, err := expectArgLength(arguments, 2); err != nil {
			return nil, err
		}

		start, err := getArg(e, JoInt, arguments[0])

		if err != nil {
			return nil, err
		}

		startLit, _ := start.(LiteralData)
		startInt := int(startLit.IntVal)

		if startInt < 0 || startInt > len(a.data)-1 {
			return nil, ErrIndexOutofBound
		}

		end, err := getArg(e, JoInt, arguments[1])

		if err != nil {
			return nil, err
		}

		endLit, _ := end.(LiteralData)
		endInt := int(endLit.IntVal)

		if endInt < 0 || endInt > len(a.data) {
			return nil, ErrIndexOutofBound
		}

		return NewArray(a.data[startInt:endInt]), nil
	case "type":
		if _, err := expectArgLength(arguments, 0); err != nil {
			return nil, err
		}
		return StringLiteral(a.Type()), nil
	}
	return nil, ErrNoMethod(name, a.Type())
}
