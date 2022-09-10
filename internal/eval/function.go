package eval

import (
	"fmt"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

func (e *Evaluator) functionDecl(node Node.Node) (EnvironmentData, error) {
	functionDecl := node.(*Node.FunctionDeclStatement)

	functionName := functionDecl.Identifier.(*Node.Identifier)

	if _, err := e.environment.Get(functionName.Value); err == nil {
		return nil, e.NewError(functionName.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", functionName.Value))
	}
	e.environment.Define(functionName.Value, NewCallableFunction(*functionDecl, e.environment, nil))
	return nil, nil
}

func (e *Evaluator) functionCall(node Node.Node) (EnvironmentData, error) {
	functionCall := node.(*Node.FunctionCall)

	functionName, _ := functionCall.Identifier.(*Node.Identifier)

	var function EnvironmentData
	switch functionCall.Identifier.NodeName() {
	case Node.IDENTIFIER:
		if functionName.Value == "print" {
			output := ""
			for i, arg := range functionCall.Arguments {
				exp, err := e.EvalExpression(arg)

				if err != nil {
					return nil, err
				}

				if i > 0 {
					output += " "
				}

				if exp == nil {
					output += "null"
				} else {
					output += exp.GetString()
				}
			}
			stdio.Io.Println(output)
			return nil, nil
		} else if functionName.Value == "input" {
			if len(functionCall.Arguments) != 1 {
				e.NewError(functionName.Token, JoError.DefaultError, "must have 1 argument.")
				return nil, e.NewError(functionName.Token, JoError.DefaultError, "must have 1 argument.")
			}
			arg1 := functionCall.Arguments[0]
			arg, err := e.EvalExpression(arg1)

			if err != nil {
				return nil, err
			}
			argLiteral := arg.(LiteralData)

			stdio.Io.Print(argLiteral.GetString())

			text := stdio.Io.Input()
			return StringLiteral(text), nil
		}

		fun, err := e.environment.Get(functionName.Value)
		if err != nil {

			return nil, e.NewError(functionName.Token, JoError.DefaultError, fmt.Sprintf("unknown function ` %s `", functionName.Value))
		}
		function = fun
	case Node.GET_EXPR:
		_structMethod, err := e._get(functionCall.Identifier)
		if err != nil {
			return nil, err
		}

		function = _structMethod
	case Node.FUNCTION_CALL:
		fun, err := e.functionCall(functionCall.Identifier)
		if err != nil {
			return nil, err
		}
		function = fun
	default:
		e.NewError(functionName.Token, JoError.DefaultError, "cannot call struct data")
		return nil, e.NewError(functionName.Token, JoError.DefaultError, "cannot call struct data")
	}

	callableFunction, ok := function.(*CallableFunction)

	if ok {
		// if callableFunction.parent != nil {
		// TODO Struct attributes
		// fmt.Println("METHOD")
		// &callableFunction.parent.env.
		// callableFunction.
		// }

		a, err := callableFunction.Call(e, node, functionCall.Arguments)
		return a, err
	}

	structDecl, ok := function.(*StructDataDecl)

	if ok {
		data := NewStructData(*structDecl)

		return data, nil
	}

	_, ok = function.(*StructData)

	if ok {
		return nil, e.NewError(e.NewTokenFromLine(functionCall.GetLine()), JoError.DefaultError, "cannot call struct data")
	}

	return nil, e.NewError(e.NewTokenFromLine(functionCall.GetLine()), JoError.DefaultError, "not a function")
}

func (e *Evaluator) Return(node Node.Node) (EnvironmentData, error) {
	returnStmt := node.(*Node.ReturnStatement)

	if returnStmt.Expression == nil {
		return StringLiteral("null"), nil
	}

	val, err := e.EvalExpression(returnStmt.Expression)

	if err != nil {
		return val, err
	}

	return val, nil
}
