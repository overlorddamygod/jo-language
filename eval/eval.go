package eval

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"

	L "github.com/overlorddamygod/jo/lexer"
	"github.com/overlorddamygod/jo/parser"
)

type Evaluator struct {
	lexer       *L.Lexer
	node        []parser.Node
	global      *Environment
	environment *Environment
	current     parser.Node
	// variables   map[string]parser.LiteralValue
}

func NewEvaluator(lexer *L.Lexer, node []parser.Node) *Evaluator {
	env := NewEnvironment()
	return &Evaluator{lexer: lexer, node: node, global: env, environment: env}
}

func NewEvaluatorWithParent(e *Evaluator, parent *Environment) *Evaluator {
	// env := NewEnvironment()
	env := NewEnvironmentWithParent(parent)
	return &Evaluator{lexer: e.lexer, node: e.node, global: env, environment: env}
}

func (e *Evaluator) Eval() (EnvironmentData, error) {
	return e.EvalStatements(e.node)
}

func (e *Evaluator) EvalStatements(statements []parser.Node) (EnvironmentData, error) {

	// for _, s := range statements {
	// 	s.Print()
	// }
	for _, s := range statements {
		data, err := e.EvalStatement(s)
		// fmt.Println("EVALSTATEMENT", s, data, err, err != nil)

		if err != nil {
			return nil, err
		}

		if s.NodeName() == "ReturnStatement" {
			// fmt.Println()
			return data, nil
		}

		if s.NodeName() == "IF" || s.NodeName() == "FOR" || s.NodeName() == "ReturnStatement" {
			if data != nil {
				return data, nil
			}
			// return data, nil
			// fmt.Println("IFFFF", s, data, nil)
			// return data, nil
		}
	}
	return nil, nil
}

func (e *Evaluator) EvalStatement(node parser.Node) (EnvironmentData, error) {
	// e.Eval()
	// fmt.Println("NODE", node.NodeName())
	switch node.NodeName() {
	case "VarDecl":
		varDecl := node.(*parser.VarDeclStatement)

		id := varDecl.Identifier.(*parser.Identifier)
		exp, err := e.EvalExpression(*varDecl.Expression)

		if err != nil {
			return nil, err
		}

		if _, err := e.environment.Get(id.Value); err == nil {
			return nil, L.NewJoError(e.lexer, id.Token, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
		} else {
			e.environment.Define(id.Value, exp)
		}
	case "ASSIGNMENT":
		assignment := node.(*parser.AssignmentStatement)
		// fmt.Println("ASSIGNMENT", assignment.Identifier)

		id := assignment.Identifier.(*parser.Identifier)
		exp, err := e.EvalExpression(assignment.Expression)

		if err != nil {
			return nil, err
		}

		err = e.environment.Assign(id.Value, exp)

		if err != nil {
			return nil, L.NewJoError(e.lexer, id.Token, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
		}
	case "FunctionCall":
		return e.functionCall(node)
	case "IF":
		// fmt.Println("IF Start")
		ifStatement := node.(*parser.IfStatement)
		e.begin()

		if ifStatement.HasIfs() {
			for _, block := range ifStatement.IfBlocks {
				literalData, err := e.EvalExpression(block.Condition)

				if err != nil {
					return nil, err
				}

				literal := literalData.(LiteralData)

				if literal.GetBoolean() {
					data, err := e.EvalStatements(block.Block.Nodes)
					if err != nil {
						e.end()
						return data, err
					}
					e.end()
					return data, err
				}
			}
		}

		if !ifStatement.HasElse() {
			e.end()
			return nil, nil
		}
		data, err := e.EvalStatements(ifStatement.ElseBlock.Nodes)

		if err != nil {
			e.end()
			return data, err
		}
		e.end()
		return data, err
	case "FOR":
		// fmt.Println("FOR CALL")
		forStatement := node.(*parser.ForStatement)
		prev := e.current
		e.current = forStatement
		e.begin()

		data, err := e.EvalStatement(forStatement.Initial)

		if err != nil {
			return data, nil
		}

		for {
			conditionData, err := e.EvalExpression(forStatement.Condition)

			if err != nil {
				e.end()
				return nil, err
			}

			condition := conditionData.(LiteralData)

			// fmt.Println("FOR condition", err)

			if !condition.GetBoolean() {
				if prev == nil || prev.NodeName() != "FOR" {
					e.current = nil
				}
				e.end()
				break
			}

			data, err = e.EvalStatements(forStatement.Block.Nodes)

			// fmt.Println("FOR CALL block", data, err)

			if err != nil {
				if err.Error() == "Statement:Break" {
					if prev == nil || prev.NodeName() != "FOR" {
						e.current = nil
					}

					e.end()
					break
				}
				if err.Error() == "Statement:Continue" {
					if prev == nil || prev.NodeName() != "FOR" {
						e.current = nil
					}

					_, err = e.EvalStatement(forStatement.Expression)

					if err != nil {
						return nil, err
					}

					continue
				}
				if prev == nil || prev.NodeName() != "FOR" {
					e.current = nil
					// e.end()
				}
				// e.end()

				return nil, err
			}

			if data != nil {
				// e.end()
				return data, nil
			}

			_, err = e.EvalStatement(forStatement.Expression)

			if err != nil {
				return nil, err
			}
		}
	case "FunctionDecl":
		functionDecl := node.(*parser.FunctionDeclStatement)

		functionName := functionDecl.Identifier.(*parser.Identifier)
		// fmt.Println("GLOBASTART-----")
		// e.global.Print()
		// fmt.Println("GLOBALEND------")

		e.environment.Define(functionName.Value, NewCallableFunction(*functionDecl, e.environment))
		// fmt.Println("ENVSTART-----")
		// e.environment.Print()
		// fmt.Println("ENVEND------")
		// e.environment.Print()
	case "ReturnStatement":
		returnStmt := node.(*parser.ReturnStatement)

		if returnStmt.Expression == nil {
			return StringLiteral("null"), nil
		}

		val, err := e.EvalExpression(returnStmt.Expression)

		if err != nil {
			return val, err
		}

		// return val, errors.New("Statement:Return")
		return val, nil
	case "BreakStatement":
		if e.current != nil && e.current.NodeName() == "FOR" {
			return nil, errors.New("Statement:Break")
		}
		// fmt.Println("HERE", e.current)
		//  := node.(*parser.BreakStatement)
		return nil, nil
	case "ContinueStatement":
		if e.current != nil && e.current.NodeName() == "FOR" {
			return nil, errors.New("Statement:Continue")
		}
		// fmt.Println("HERE", e.current)
		//  := node.(*parser.BreakStatement)
		return nil, nil
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
		for i, arg := range functionCall.Arguments {
			exp, err := e.EvalExpression(arg)

			if err != nil {
				return nil, err
			}

			if i > 0 {
				output += " "
			}
			expressionVal, ok := exp.(LiteralData)

			if !ok {
				output += "null"
			} else {
				output += expressionVal.Value
			}
		}
		fmt.Println(output)
		return nil, nil
	} else if functionName.Value == "input" {
		if len(functionCall.Arguments) != 1 {
			return nil, L.NewJoError(e.lexer, functionName.Token, "must have 1 argument.")
		}
		arg1 := functionCall.Arguments[0]
		arg, err := e.EvalExpression(arg1)

		if err != nil {
			return nil, err
		}
		argLiteral := arg.(LiteralData)

		fmt.Print(argLiteral.GetString())
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		return StringLiteral(scanner.Text()), nil
	}
	// fmt.Println("FUNC START", functionName.Value)

	function, err := e.environment.Get(functionName.Value)
	if err != nil {
		return nil, L.NewJoError(e.lexer, functionName.Token, fmt.Sprintf("unknown function ` %s `", functionName.Value))
	}

	callableFunction := function.(*CallableFunction)

	a, err := callableFunction.Call(e, functionCall.Arguments)

	// fmt.Println("FUNC END", functionName.Value, a, err)

	return a, err
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
