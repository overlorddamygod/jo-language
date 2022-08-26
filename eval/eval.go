package eval

import (
	"fmt"
	"math"

	L "github.com/overlorddamygod/jo/lexer"
	"github.com/overlorddamygod/jo/parser"
)

type Evaluator struct {
	node      []parser.Node
	variables map[string]parser.LiteralValue
}

func NewEvaluator(node []parser.Node) *Evaluator {
	return &Evaluator{node: node, variables: make(map[string]parser.LiteralValue)}
}

func (e *Evaluator) Eval() {
	e.EvalStatements(e.node)
}

func (e *Evaluator) EvalStatements(statements []parser.Node) {
	for _, s := range statements {
		e.EvalStatement(s)
	}
}

func (e *Evaluator) EvalStatement(node parser.Node) {
	// e.Eval()
	switch node.NodeName() {
	case "ASSIGNMENT":
		assignment := node.(*parser.AssignmentStatement)
		// fmt.Println("ASSIGNMENT", assignment.Identifier)

		id := assignment.Identifier.(*parser.Identifier)
		exp := e.EvalExpression(assignment.Expression)
		e.variables[id.Value] = exp
		// if binaryExpression.Op == "ASSIGNMENT" {
		// fmt.Println(e.EvalExpression(assignment.Expression))
		// }
	case "FUNCTION_CALL":
		functionCall := node.(*parser.FunctionCallStatement)

		functionName := functionCall.Identifier.(*parser.Identifier)

		if functionName.Value == "print" {
			// fmt.Println("HERE")
			exp := e.EvalExpression(functionCall.Expression)
			fmt.Println(exp.Value)
		}
	case "IF":
		ifStatement := node.(*parser.IfStatement)

		literal := e.EvalExpression(ifStatement.Condition)

		if literal.GetBoolean() {
			e.EvalStatements(ifStatement.IfBlock)
		} else {
			e.EvalStatements(ifStatement.ElseBlock)
		}
	case "FOR":
		forStatement := node.(*parser.ForStatement)

		e.EvalStatement(forStatement.Initial)

		for {
			condition := e.EvalExpression(forStatement.Condition)

			if !condition.GetBoolean() {
				break
			}

			e.EvalStatements(forStatement.Block)

			e.EvalStatement(forStatement.Expression)
		}
	default:
		fmt.Println("UNKNOWN STATEMENT", node.NodeName())
	}
}

func (e *Evaluator) EvalExpression(node parser.Node) parser.LiteralValue {
	switch node.NodeName() {
	case "BinaryExpression":
		binaryExpression := node.(*parser.BinaryExpression)

		left := e.EvalExpression(binaryExpression.Left)
		right := e.EvalExpression(binaryExpression.Right)

		if left.IsNumber() {
			switch binaryExpression.Op {
			case L.PLUS:
				return parser.NumberLiteral(left.GetNumber() + right.GetNumber())
			case L.MINUS:
				return parser.NumberLiteral(left.GetNumber() - right.GetNumber())
			case L.SLASH:
				return parser.NumberLiteral(left.GetNumber() / right.GetNumber())
			case L.ASTERISK:
				return parser.NumberLiteral(left.GetNumber() * right.GetNumber())
			case L.PERCENT:
				return parser.NumberLiteral(math.Mod(left.GetNumber(), right.GetNumber()))
			case L.EQ:
				return parser.BooleanLiteral(left.GetNumber() == right.GetNumber())
			case L.NOT_EQ:
				return parser.BooleanLiteral(left.GetNumber() != right.GetNumber())
			case L.GT:
				return parser.BooleanLiteral(left.GetNumber() > right.GetNumber())
			case L.GT_EQ:
				return parser.BooleanLiteral(left.GetNumber() >= right.GetNumber())
			case L.LT:
				return parser.BooleanLiteral(left.GetNumber() < right.GetNumber())
			case L.LT_EQ:
				return parser.BooleanLiteral(left.GetNumber() <= right.GetNumber())
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
			return parser.BooleanLiteral(val)
		}

		if left.IsString() {
			switch binaryExpression.Op {
			case L.PLUS:
				return parser.StringLiteral(left.GetString() + right.GetString())
			case L.EQ:
				return parser.BooleanLiteral(left.GetString() == right.GetString())
			case L.NOT_EQ:
				return parser.BooleanLiteral(left.GetString() != right.GetString())
			}
		}

	case "LiteralValue":
		literal := node.(*parser.LiteralValue)
		return *literal
	case "Identifier":
		variable := node.(*parser.Identifier)
		return e.variables[variable.Value]
	default:
		fmt.Println("Unknown Nodename")
	}
	return parser.NumberLiteral(0)
}
