package eval

import (
	"errors"
	"fmt"
)

type EnvironmentData interface {
	Type() string
	GetString() string
}

type EnvironmentDataValue EnvironmentData
type EnvironmentDataMap map[string]EnvironmentDataValue

type Environment struct {
	data   EnvironmentDataMap
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		data:   make(EnvironmentDataMap, 0),
		parent: nil,
	}
}

func NewEnvironmentWithParent(env *Environment) *Environment {
	return &Environment{
		data:   make(EnvironmentDataMap, 0),
		parent: env,
	}
}

func (env *Environment) Define(key string, value EnvironmentDataValue) {
	env.data[key] = value
}

func (env *Environment) Print() {
	fmt.Println("KEYS")
	if env == nil {
		return
	}
	for key, val := range env.data {
		// fmt.Println(key)
		if val.Type() == Function {
			f := val.(*CallableFunction)

			println("FUNC", key, f.Type())
		} else if val.Type() == string(Literal) {
			lit := val.(LiteralData)
			println("VAL", key, lit.Value)
		} else if val.Type() == Struct {
			s := val.(*StructData)
			println("Struct", key, s)
		} else {
			println("Struct Decl", key, val)
		}
		// if val.Type() == "LiteralData" {
		// } else {
		// println("fun", key, val.Type())
		// }
	}
}

func (env *Environment) Get(key string) (EnvironmentDataValue, error) {
	value, present := env.data[key]

	if !present {
		if env.parent == nil {
			return nil, errors.New("key not defined in the environment")
		}

		return env.parent.Get(key)
	}
	return value, nil
}

func (env *Environment) GetOne(key string) (EnvironmentDataValue, error) {
	value, present := env.data[key]

	if !present {
		return nil, errors.New("key not defined in the environment")
	}
	return value, nil
}

func (env *Environment) Assign(key string, value EnvironmentDataValue) error {
	_, present := env.data[key]

	if !present {
		if env.parent == nil {
			return errors.New("key not defined in the environment")
		}
		return env.parent.Assign(key, value)
	}

	env.data[key] = value
	return nil
}
