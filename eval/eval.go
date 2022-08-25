package eval

import (
	"fmt"
	"strconv"

	L "github.com/overlorddamygod/lexer/lexer"
	"github.com/overlorddamygod/lexer/parser"
)

type Evaluator struct {
	node      []parser.Node
	variables map[string]float64
}

func NewEvaluator(node []parser.Node) *Evaluator {
	return &Evaluator{node: node, variables: make(map[string]float64)}
}

func (e *Evaluator) Eval() {

	for _, s := range e.node {
		e.EvalStatement(s)
	}
	// e.EvalStatement(e.node)
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
			fmt.Println(exp)
		}
	}
}

func (e *Evaluator) EvalExpression(node parser.Node) float64 {
	switch node.NodeName() {
	case "BinaryExpression":
		binaryExpression := node.(*parser.BinaryExpression)

		left := e.EvalExpression(binaryExpression.Left)
		right := e.EvalExpression(binaryExpression.Right)
		switch binaryExpression.Op {
		case L.PLUS:
			return left + right
		case L.MINUS:
			return left - right
		case L.SLASH:
			return left / right
		case L.ASTERISK:
			return left * right
		}

		// case L.PLUS:
		// 	return e.EvalExpression(binaryExpression.Left) + e.EvalExpression(binaryExpression.Right)
		// }
	case "LiteralValue":
		literal := node.(*parser.LiteralValue)
		val, _ := strconv.ParseFloat(literal.Value, 32)
		return val
	case "Identifier":
		variable := node.(*parser.Identifier)
		return e.variables[variable.Value]
	default:
		fmt.Println("Unknown Nodename")
	}
	return 0
}
