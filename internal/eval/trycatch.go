package eval

import (
	"errors"

	joerror "github.com/overlorddamygod/jo/pkg/error"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) tryCatch(node Node.Node) (EnvironmentData, error) {
	tryCatchStmt := node.(*Node.TryCatchStatement)

	prevTryCatchScope := e.TryCatchScope
	e.TryCatchScope = true

	defer func() {
		e.TryCatchScope = prevTryCatchScope
	}()
	// functionCa

	e.begin()
	data, err := e.EvalStatements(tryCatchStmt.Try.Nodes)
	e.end()
	if err != nil {
		if errors.Is(err, ErrThrow) {
			e.begin()
			e.environment.DefineOne(tryCatchStmt.CatchVar.Value, data)
			data, err = e.EvalStatements(tryCatchStmt.Catch.Nodes)
			e.end()

			if err != nil {
				return data, err
			}

			if data != nil {
				return data, nil
			}
		}
		return nil, err
	}
	return data, nil
}

func (e *Evaluator) throw(node Node.Node) (EnvironmentData, error) {

	throwStmt := node.(*Node.ThrowStatement)

	if !e.TryCatchScope {
		return nil, e.NewError(e.NewTokenFromLine(throwStmt.GetLine()), joerror.DefaultError, "not inside a try catch block")
	}
	// if ThrowStmt.Expression == nil {
	// 	return NullLiteral(), nil
	// }

	val, err := e.EvalExpression(throwStmt.Expression)

	if err != nil {
		return val, err
	}

	return val, ErrThrow
}
