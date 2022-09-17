package eval

import (
	"errors"

	L "github.com/overlorddamygod/jo/pkg/lexer"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) IfElse(node Node.Node) (EnvironmentData, error) {
	ifStatement := node.(*Node.IfStatement)

	if ifStatement.HasIfs() {
		for _, block := range ifStatement.IfBlocks {
			data, err := e.EvalExpression(block.Condition)

			if err != nil {
				return nil, err
			}

			if data.GetBoolean() {
				e.begin()
				data, err := e.EvalStatements(block.Block.Nodes)
				e.end()
				if err != nil {
					return data, err
				}
				return data, err
			}
		}
	}

	if !ifStatement.HasElse() {
		return nil, nil
	}

	e.begin()
	data, err := e.EvalStatements(ifStatement.ElseBlock.Nodes)
	e.end()
	if err != nil {
		return data, err
	}
	return data, err
}

func (e *Evaluator) match(testVal EnvironmentData, values []Node.Node) (bool, error) {
	for _, val := range values {
		evaluatedVal, err := e.EvalExpression(val)
		if err != nil {
			return false, err
		}

		boolVal, err := e.BinaryOp(testVal, L.EQ, evaluatedVal)

		if err != nil {
			return false, err
		}

		if boolVal.Type() != JoBoolean {
			return false, ErrUnexpected
		}

		boolLiteral := boolVal.(LiteralData)

		if boolLiteral.GetBoolean() {
			return true, nil
		}
	}
	return false, nil
}

func (e *Evaluator) Switch(node Node.Node) (EnvironmentData, error) {
	switchStmt := node.(*Node.SwitchStatement)
	prev := e.current
	e.current = switchStmt
	defer func() {
		e.current = prev
	}()
	testValue, err := e.EvalExpression(switchStmt.Test)

	if err != nil {
		return nil, err
	}

	caseMatched := false
	for _, _case := range switchStmt.Cases {
		match, err := e.match(testValue, _case.Values)

		if err != nil {
			return nil, err
		}

		if match {
			caseMatched = true
			e.begin()
			data, err := e.EvalStatements(_case.Block.Nodes)
			e.end()

			if err != nil {
				if errors.Is(err, ErrBreak) {
					return nil, nil
				}
				return nil, err
			}

			if data != nil {
				return data, nil
			}
		}
	}

	// default

	if !caseMatched {
		e.begin()
		data, err := e.EvalStatements(switchStmt.Default.Nodes)
		e.end()
		if err != nil {
			if errors.Is(err, ErrBreak) {
				return nil, nil
			}
			return nil, err
		}

		if data != nil {
			return data, nil
		}
	}
	return nil, nil
}
