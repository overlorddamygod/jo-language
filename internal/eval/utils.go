package eval

import (
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func expectArgLength(arg []Node.Node, length int) (EnvironmentData, error) {
	if length == -1 {
		return nil, nil
	}
	if len(arg) != length {
		return nil, ErrArgumentLength(length)
	}
	return nil, nil
}

func getArg(e *Evaluator, Type string, node Node.Node) (EnvironmentData, error) {
	evaluatedData, err := e.EvalExpression(node)

	if err != nil {
		return nil, err
	}

	if evaluatedData.Type() != Type {
		return nil, ErrArgumentType(Type)
	}

	return evaluatedData, nil
}
