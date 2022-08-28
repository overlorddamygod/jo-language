package eval

import "errors"

type EnvironmentData interface {
	Type() string
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
