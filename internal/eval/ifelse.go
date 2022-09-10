package eval

import Node "github.com/overlorddamygod/jo/pkg/parser/node"

func (e *Evaluator) IfElse(node Node.Node) (EnvironmentData, error) {
	ifStatement := node.(*Node.IfStatement)

	if ifStatement.HasIfs() {
		for _, block := range ifStatement.IfBlocks {
			literalData, err := e.EvalExpression(block.Condition)

			if err != nil {
				return nil, err
			}

			literal := literalData.(LiteralData)

			if literal.GetBoolean() {
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
