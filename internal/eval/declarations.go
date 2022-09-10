package eval

import (
	"fmt"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) structDecl(node Node.Node) (EnvironmentData, error) {
	structD := node.(*Node.StructDeclStatement)

	id := structD.Identifier.(*Node.Identifier)

	if _, err := e.environment.Get(id.Value); err == nil {
		return nil, e.NewError(id.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
	}
	e.environment.Define(id.Value, NewStructDataDecl(*structD, e.environment))
	return nil, nil
}

func (e *Evaluator) varDecl(node Node.Node) (EnvironmentData, error) {

	varDecl := node.(*Node.VarDeclStatement)

	id := varDecl.Identifier.(*Node.Identifier)
	exp, err := e.EvalExpression(*varDecl.Expression)

	if err != nil {
		return nil, err
	}

	if _, err := e.environment.GetOne(id.Value); err == nil {
		return nil, e.NewError(id.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
	}
	e.environment.Define(id.Value, exp)
	return nil, nil
}
