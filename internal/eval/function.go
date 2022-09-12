package eval

import (
	"errors"
	"fmt"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
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
	// functionCall.Print()
	// fmt.Println(functionCall.Identifier, functionCall.Arguments)

	functionName, _ := functionCall.Identifier.(*Node.Identifier)

	var function EnvironmentData
	switch functionCall.Identifier.NodeName() {
	case Node.IDENTIFIER:
		fun, err := e.environment.Get(functionName.Value)
		if err != nil {

			return nil, e.NewError(functionName.Token, JoError.DefaultError, fmt.Sprintf("unknown function ` %s `", functionName.Value))
		}
		function = fun
	case Node.GET_EXPR:
		getexpr := functionCall.Identifier.(*Node.GetExpr)

		name, _ := getexpr.Identifier.(*Node.Identifier)

		left, err := e.EvalExpression(getexpr.Expr)
		// fmt.Println("ZZZ", left)
		if err != nil {
			return nil, err
		}

		// Todo merge all types
		data, ok := left.(LiteralData)
		if ok {
			name, _ := getexpr.Identifier.(*Node.Identifier)
			return data.Call(e, name.Value, functionCall.Arguments)
		}

		arrayData, ok := left.(*Array)
		if ok {
			name, _ := getexpr.Identifier.(*Node.Identifier)
			returnData, err := arrayData.Call(e, name.Value, functionCall.Arguments)

			if err != nil {
				return nil, e.NewError(name.Token, JoError.DefaultError, err.Error())
			}
			return returnData, nil
		}

		structData, ok := left.(*StructData)
		if ok {
			returnData, err := structData.Call(e, name.Value, functionCall.Arguments)

			if err != nil {
				return nil, e.NewError(name.Token, JoError.DefaultError, err.Error())
			}
			return returnData, nil
		}

		// node.Print()
		return nil, errors.New("unknown")
		// function = val
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
		name := functionCall.Identifier.(*Node.Identifier)
		a, err := callableFunction.Call(e, name.Value, functionCall.Arguments)

		if err != nil {
			return nil, e.NewError(functionName.Token, JoError.DefaultError, err.Error())
		}
		return a, err
	}

	c, ok := function.(*CallableFunc)

	if ok {
		name := functionCall.Identifier.(*Node.Identifier)
		if c.Arity == -1 || len(functionCall.Arguments) == c.Arity {
			return c.Call(e, name.Value, functionCall.Arguments)
		}
		return nil, e.NewError(functionName.Token, JoError.DefaultError, "Arguments length doesnot match")
	}

	structDecl, ok := function.(*StructDataDecl)

	if ok {
		return structDecl.Initialize(e, functionCall.Arguments)
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
