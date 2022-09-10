package eval

import (
	"errors"

	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

var (
	ErrBreak    = errors.New("Statement:Break")
	ErrContinue = errors.New("Statement:Continue")
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

		condition := conditionData.(LiteralData)

		if !condition.GetBoolean() {
			e.end()
			break
		}

		_, err = e.EvalStatements(forStatement.Block.Nodes)
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
			return nil, err
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
		condition := conditionData.(LiteralData)

		if !condition.GetBoolean() {
			e.end()
			break
		}

		_, err = e.EvalStatements(whileStatement.Block.Nodes)
		e.end()

		if err != nil {
			if errors.Is(err, ErrBreak) {
				break
			}
			if errors.Is(err, ErrContinue) {
				continue
			}
			return nil, err
		}
	}
	e.current = prev
	return nil, nil
}
