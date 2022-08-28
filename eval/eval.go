package eval

import (
	"errors"
	"fmt"
	"math"

	L "github.com/overlorddamygod/jo/lexer"
	"github.com/overlorddamygod/jo/parser"
)

type Evaluator struct {
	lexer       *L.Lexer
	node        []parser.Node
	global      *Environment
	environment *Environment
	// variables   map[string]parser.LiteralValue
}

func NewEvaluator(lexer *L.Lexer, node []parser.Node) *Evaluator {
	env := NewEnvironment()
	return &Evaluator{lexer: lexer, node: node, global: env, environment: env}
}

func (e *Evaluator) Eval() (EnvironmentData, error) {
	return e.EvalStatements(e.node)
}

func (e *Evaluator) EvalStatements(statements []parser.Node) (EnvironmentData, error) {
	for _, s := range statements {
		data, err := e.EvalStatement(s)
		// fmt.Println("EVALSTATEMENT", s, data, err, err != nil)
		if err != nil {
			// fmt.Println("HERE", err)
			return nil, err
		}
		if data != nil {
			return data, nil
		}
	}
	return nil, nil
}

func (e *Evaluator) EvalStatement(node parser.Node) (EnvironmentData, error) {
	// e.Eval()
	switch node.NodeName() {
	case "ASSIGNMENT":
		assignment := node.(*parser.AssignmentStatement)
		// fmt.Println("ASSIGNMENT", assignment.Identifier)

		id := assignment.Identifier.(*parser.Identifier)
		exp, err := e.EvalExpression(assignment.Expression)

		if err != nil {
			return nil, err
		}
		e.environment.Define(id.Value, exp)
	case "FunctionCall":
		return e.functionCall(node)
	case "IF":
		// fmt.Println("IF Start")

		ifStatement := node.(*parser.IfStatement)

		literalData, err := e.EvalExpression(ifStatement.Condition)

		// fmt.Println("IF EXP", literalData, err)

		if err != nil {
			return nil, err
		}

		// fmt.Println("IF EXPzzzzzz", literalData, err)

		literal := literalData.(LiteralData)
		// fmt.Println("IF", err)

		e.begin()
		if literal.GetBoolean() {
			data, err := e.EvalStatements(ifStatement.IfBlock)
			// fmt.Println("IF true", data, err)
			if err != nil {
				return nil, err
			}

			e.end()
			return data, err
		} else {
			data, err := e.EvalStatements(ifStatement.ElseBlock)

			if err != nil {
				return nil, err
			}
			e.end()
			return data, err
		}
	case "FOR":
		// fmt.Println("FOR CALL")
		forStatement := node.(*parser.ForStatement)
		e.begin()

		data, err := e.EvalStatement(forStatement.Initial)

		if err != nil {
			return data, nil
		}
		// fmt.Println("FOR Init", err)

		for {
			conditionData, err := e.EvalExpression(forStatement.Condition)

			if err != nil {
				e.end()
				return nil, err
			}

			condition := conditionData.(LiteralData)

			// fmt.Println("FOR condition", err)

			if !condition.GetBoolean() {
				e.end()
				break
			}

			_, err = e.EvalStatements(forStatement.Block)
			// fmt.Println("FOR CALL block", err)

			if err != nil {
				return nil, err
			}

			_, err = e.EvalStatement(forStatement.Expression)
			// fmt.Println("FOR CALL right", err)

			if err != nil {
				return nil, err
			}
		}
	case "FunctionDecl":
		functionDecl := node.(*parser.FunctionDeclStatement)

		functionName := functionDecl.Identifier.(*parser.Identifier)

		e.environment.Define(functionName.Value, NewCallableFunction(*functionDecl))
	case "ReturnStatement":
		returnStmt := node.(*parser.ReturnStatement)

		return e.EvalExpression(returnStmt.Expression)
	default:
		return nil, fmt.Errorf("unknown statement %s", node.NodeName())
	}
	return nil, nil
}

func (e *Evaluator) EvalExpression(node parser.Node) (EnvironmentData, error) {
	switch node.NodeName() {
	case "BinaryExpression":
		binaryExpression := node.(*parser.BinaryExpression)

		leftData, err := e.EvalExpression(binaryExpression.Left)

		if err != nil {
			return nil, err
		}

		// if left.Type() == "LiteralData" {
		// 	left = left.(LiteralData)
		// }
		rightData, err := e.EvalExpression(binaryExpression.Right)

		if err != nil {
			return nil, err
		}

		left := leftData.(LiteralData)
		right := rightData.(LiteralData)

		// if right.Type() == "LiteralData" {
		// 	right = right.(LiteralData)
		// }

		if left.IsNumber() {
			switch binaryExpression.Op {
			case L.PLUS:
				return NumberLiteral(left.GetNumber() + right.GetNumber()), nil
			case L.MINUS:
				return NumberLiteral(left.GetNumber() - right.GetNumber()), nil
			case L.SLASH:
				return NumberLiteral(left.GetNumber() / right.GetNumber()), nil
			case L.ASTERISK:
				return NumberLiteral(left.GetNumber() * right.GetNumber()), nil
			case L.PERCENT:
				return NumberLiteral(math.Mod(left.GetNumber(), right.GetNumber())), nil
			case L.EQ:
				return BooleanLiteral(left.GetNumber() == right.GetNumber()), nil
			case L.NOT_EQ:
				return BooleanLiteral(left.GetNumber() != right.GetNumber()), nil
			case L.GT:
				return BooleanLiteral(left.GetNumber() > right.GetNumber()), nil
			case L.GT_EQ:
				return BooleanLiteral(left.GetNumber() >= right.GetNumber()), nil
			case L.LT:
				return BooleanLiteral(left.GetNumber() < right.GetNumber()), nil
			case L.LT_EQ:
				return BooleanLiteral(left.GetNumber() <= right.GetNumber()), nil
			}
		}

		if left.IsBoolean() {
			var val bool
			switch binaryExpression.Op {
			case L.EQ:
				val = left.GetBoolean() == right.GetBoolean()
			case L.NOT_EQ:
				val = left.GetBoolean() != right.GetBoolean()
			case L.AND:
				val = left.GetBoolean() && right.GetBoolean()
			case L.OR:
				val = left.GetBoolean() || right.GetBoolean()
			}
			return BooleanLiteral(val), nil
		}

		if left.IsString() {
			switch binaryExpression.Op {
			case L.PLUS:
				return StringLiteral(left.GetString() + right.GetString()), nil
			case L.EQ:
				return BooleanLiteral(left.GetString() == right.GetString()), nil
			case L.NOT_EQ:
				return BooleanLiteral(left.GetString() != right.GetString()), nil
			}
		}

	case "LiteralValue":
		literal := node.(*parser.LiteralValue)
		return LiteralDataFromParserLiteral(*literal), nil
	case "UnaryExpression":
		unary := node.(*parser.UnaryExpression)

		if unary.Op == L.BANG {
			data, err := e.EvalExpression(unary.Identifier)
			if err != nil {
				return nil, err
			}

			value := data.(LiteralData)

			return BooleanLiteral(!value.GetBoolean()), nil
		}
		return nil, L.NewJoError(e.lexer, unary.Token, "Unknown operator "+unary.Op)
	case "Identifier":
		variable := node.(*parser.Identifier)
		val, err := e.environment.Get(variable.Value)

		if err != nil {
			// fmt.Println(err)
			return nil, L.NewJoError(e.lexer, variable.Token, fmt.Sprintf("Variable ` %s ` not defined in this scope", variable.Value))
		}
		return val, nil
	case "FunctionCall":
		return e.functionCall(node)
	default:
		return nil, errors.New("unknown nodename")
	}
	return nil, nil
}

func (e *Evaluator) functionCall(node parser.Node) (EnvironmentData, error) {
	functionCall := node.(*parser.FunctionCall)

	functionName := functionCall.Identifier.(*parser.Identifier)

	if functionName.Value == "print" {
		output := ""
		for _, arg := range functionCall.Arguments {
			exp, err := e.EvalExpression(arg)

			if err != nil {
				return nil, err
			}

			expressionVal := exp.(LiteralData)
			output += " " + expressionVal.Value
		}
		fmt.Println(output)
		return nil, nil
	}

	function, err := e.environment.Get(functionName.Value)
	if err != nil {
		return nil, L.NewJoError(e.lexer, functionName.Token, fmt.Sprintf("unknown function ` %s `", functionName.Value))
	}

	callableFunction := function.(*CallableFunction)

	return callableFunction.Exec(e, functionCall.Arguments)
}

func (e *Evaluator) begin() {
	e.environment = NewEnvironmentWithParent(e.environment)
}

func (e *Evaluator) end() {
	if e.environment.parent == nil {
		e.environment = e.global
		return
	}
	e.environment = e.environment.parent
}
