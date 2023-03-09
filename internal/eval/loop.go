package eval

import (
	"errors"

	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) For(node Node.Node) (EnvironmentData, error) {
	defer func() {
		e.end()
	}()
	forStatement := node.(*Node.ForStatement)
	prev := e.current
	e.current = forStatement
	e.begin()

	data, err := e.EvalStatement(forStatement.Initial)

	if err != nil {
		return data, nil
	}

	for {
		e.begin()
		conditionData, err := e.EvalExpression(forStatement.Condition)

		if err != nil {
			e.end()
			return nil, err
		}

		if !conditionData.GetBoolean() {
			e.end()
			break
		}

		data, err = e.EvalStatements(forStatement.Block.Nodes)
		e.end()

		if err != nil {
			if errors.Is(err, ErrBreak) {
				break
			}
			if errors.Is(err, ErrContinue) {
				_, err = e.EvalStatement(forStatement.Expression)

				if err != nil {
					return nil, err
				}

				continue
			}

			if errors.Is(err, ErrThrow) {
				return data, err
			}
			return nil, err
		}

		if data != nil {
			return data, nil
		}

		_, err = e.EvalStatement(forStatement.Expression)

		if err != nil {
			return nil, err
		}
	}
	e.current = prev
	return nil, nil
}

func (e *Evaluator) While(node Node.Node) (EnvironmentData, error) {
	whileStatement := node.(*Node.WhileStatement)
	prev := e.current
	e.current = whileStatement
	for {
		e.begin()
		conditionData, err := e.EvalExpression(whileStatement.Condition)

		if err != nil {
			e.end()
			return nil, err
		}

		if !conditionData.GetBoolean() {
			e.end()
			break
		}

		data, err := e.EvalStatements(whileStatement.Block.Nodes)
		e.end()

		if err != nil {
			if errors.Is(err, ErrBreak) {
				break
			}
			if errors.Is(err, ErrContinue) {
				continue
			}

			if errors.Is(err, ErrThrow) {
				return data, err
			}
			return nil, err
		}

		if data != nil {
			return data, nil
		}
	}
	e.current = prev
	return nil, nil
}
