package eval

import (
	"fmt"
	"math"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) BinaryOp(left EnvironmentData, op string, right EnvironmentData) (EnvironmentData, error) {
	leftData := left.(LiteralData)
	rightData := right.(LiteralData)

	// if right.Type() == "LiteralData" {
	// 	right = right.(LiteralData)
	// }

	if leftData.IsNumber() {
		switch op {
		case L.PLUS:
			return NumberLiteral(leftData.GetNumber() + rightData.GetNumber()), nil
		case L.MINUS:
			return NumberLiteral(leftData.GetNumber() - rightData.GetNumber()), nil
		case L.SLASH:
			return NumberLiteral(leftData.GetNumber() / rightData.GetNumber()), nil
		case L.ASTERISK:
			return NumberLiteral(leftData.GetNumber() * rightData.GetNumber()), nil
		case L.PERCENT:
			return NumberLiteral(math.Mod(leftData.GetNumber(), rightData.GetNumber())), nil
		case L.EQ:
			return BooleanLiteral(leftData.GetNumber() == rightData.GetNumber()), nil
		case L.NOT_EQ:
			return BooleanLiteral(leftData.GetNumber() != rightData.GetNumber()), nil
		case L.GT:
			return BooleanLiteral(leftData.GetNumber() > rightData.GetNumber()), nil
		case L.GT_EQ:
			return BooleanLiteral(leftData.GetNumber() >= rightData.GetNumber()), nil
		case L.LT:
			return BooleanLiteral(leftData.GetNumber() < rightData.GetNumber()), nil
		case L.LT_EQ:
			return BooleanLiteral(leftData.GetNumber() <= rightData.GetNumber()), nil
		case L.AND:
			return BooleanLiteral(leftData.GetBoolean() && rightData.GetBoolean()), nil
		case L.OR:
			return BooleanLiteral(leftData.GetBoolean() || rightData.GetBoolean()), nil
		}
	}

	if leftData.IsBoolean() {
		var val bool
		switch op {
		case L.EQ:
			val = leftData.GetBoolean() == rightData.GetBoolean()
		case L.NOT_EQ:
			val = leftData.GetBoolean() != rightData.GetBoolean()
		case L.AND:
			val = leftData.GetBoolean() && rightData.GetBoolean()
		case L.OR:
			val = leftData.GetBoolean() || rightData.GetBoolean()
		}
		return BooleanLiteral(val), nil
	}

	if leftData.IsString() {
		switch op {
		case L.PLUS:
			return StringLiteral(leftData.GetString() + rightData.GetString()), nil
		case L.EQ:
			return BooleanLiteral(leftData.GetString() == rightData.GetString()), nil
		case L.NOT_EQ:
			return BooleanLiteral(leftData.GetString() != rightData.GetString()), nil
		}
	}
	return NumberLiteral(2), nil
}

func (e *Evaluator) EvalExpression(node Node.Node) (EnvironmentData, error) {
	switch node.NodeName() {
	case "BinaryExpression":
		binaryExpression := node.(*Node.BinaryExpression)

		leftData, err := e.EvalExpression(binaryExpression.Left)

		if err != nil {
			return nil, err
		}

		rightData, err := e.EvalExpression(binaryExpression.Right)

		if err != nil {
			return nil, err
		}

		return e.BinaryOp(leftData, binaryExpression.Op, rightData)
	case "LiteralValue":
		literal := node.(*Node.LiteralValue)
		return LiteralDataFromParserLiteral(*literal), nil
	case "UnaryExpression":
		unary := node.(*Node.UnaryExpression)

		if unary.Op == L.BANG {
			d, err := e.EvalExpression(unary.Identifier)
			if err != nil {
				return nil, err
			}

			value := d.(LiteralData)

			return BooleanLiteral(!value.GetBoolean()), nil
		}
		return nil, e.NewError(unary.Token, JoError.DefaultError, "Unknown operator "+unary.Op)
	case "Identifier":
		return e.identifier(node)
	case "FunctionCall":
		return e.functionCall(node)
	case "GetExpr":
		return e._get(node)
	default:
		return nil, e.NewError(e.NewTokenFromLine(node.GetLine()), JoError.DefaultError, "Unknown nodename")
	}
}

func (e *Evaluator) assignment(node Node.Node) (EnvironmentData, error) {
	assignment := node.(*Node.AssignmentStatement)
	// fmt.Println("ASSIGNMENT", assignment.Identifier, assignment.Expression)

	id, ok := assignment.Identifier.(*Node.Identifier)

	if !ok {
		getExpr, _ := assignment.Identifier.(*Node.GetExpr)

		data, err := e.EvalExpression(getExpr.Expr)

		if err != nil {
			return nil, err
		}

		// fmt.Println("GET", *getExpr, getExpr.Identifier, getExpr.Expr, data, err)

		if data == nil {
			return nil, e.NewError(e.NewTokenFromLine(getExpr.GetLine()), JoError.ReferenceError, "Cannot assign to null data")
		}

		switch data.Type() {
		case Struct:
			struct_ := data.(*StructData)
			id, ok := getExpr.Identifier.(*Node.Identifier)
			// fmt.Println(id, ok)
			if ok {
				left, structGeterr := struct_.env.GetOne(id.Value)
				// if structGeterr
				// fmt.Println(left, err)

				_, ok := left.(*CallableFunction)

				if ok {
					return nil, e.NewError(id.Token, JoError.DefaultError, fmt.Sprintf("Cannot assign to method declaration ` %s `", id.Value))
				}

				// fmt.Println("STRUCT", struct_)
				exp, err := e.EvalExpression(assignment.Expression)

				if err != nil {
					return nil, err
				}
				// struct_.env.DefineOne(id.Value, exp)
				// struct_.env.Print()

				// TODO: REFACTOR THIS
				switch assignment.Op {
				case L.ASSIGN:
					err = struct_.env.DefineOne(id.Value, exp)

					if err != nil {
						return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Variable ` %s ` not defined", id.Value))
					}
				case L.PLUS, L.MINUS, L.ASTERISK, L.SLASH, L.BANG, L.PIPE, L.AND, L.OR, L.AMPERSAND, L.PERCENT:
					if structGeterr != nil {
						return nil, structGeterr
					}
					if err != nil {
						return nil, err
					}
					exp, err = e.BinaryOp(left, assignment.Op, exp)

					if err != nil {
						return nil, err
					}

					err = struct_.env.DefineOne(id.Value, exp)

					if err != nil {
						return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Variable ` %s ` not defined", id.Value))
					}
					return nil, nil
				default:
					return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Operator ` %s ` not defined", assignment.Op))
				}
			}
			return nil, nil
		default:
			return nil, e.NewError(e.NewTokenFromLine(getExpr.GetLine()), JoError.DefaultError, fmt.Sprintf("Cannot assign `%s` to the data", data.GetString()))
		}
	}

	if id.Value == "self" {
		return nil, e.NewError(id.Token, JoError.DefaultError, "Cannot assign to self keyword")
	}
	exp, err := e.EvalExpression(assignment.Expression)

	if err != nil {
		return nil, err
	}
	// fmt.Println("LOLLL", exp)
	switch assignment.Op {
	case L.ASSIGN:
		err = e.environment.Assign(id.Value, exp)

		if err != nil {
			return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Variable ` %s ` not defined", id.Value))
		}
	case L.PLUS, L.MINUS, L.ASTERISK, L.SLASH, L.BANG, L.PIPE, L.AND, L.OR, L.AMPERSAND, L.PERCENT:

		left, err := e.EvalExpression(id)

		if err != nil {
			return nil, err
		}
		exp, err = e.BinaryOp(left, assignment.Op, exp)

		if err != nil {
			return nil, err
		}

		err = e.environment.Assign(id.Value, exp)

		if err != nil {
			return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Variable ` %s ` not defined", id.Value))
		}
		return nil, nil
	default:
		return nil, e.NewError(id.Token, JoError.ReferenceError, fmt.Sprintf("Operator ` %s ` not defined", assignment.Op))
	}

	return nil, nil
}
func (e *Evaluator) identifier(node Node.Node) (EnvironmentData, error) {
	variable := node.(*Node.Identifier)
	val, err := e.environment.Get(variable.Value)

	if err != nil {
		// fmt.Println(err)
		return nil, e.NewError(variable.Token, JoError.ReferenceError, fmt.Sprintf("Variable ` %s ` not defined in this scope", variable.Value))
	}
	return val, nil
}

func (e *Evaluator) _get(node Node.Node) (EnvironmentData, error) {
	getExpr := node.(*Node.GetExpr)

	identifier := getExpr.Identifier.(*Node.Identifier)
	// fmt.Println(*&getExpr.Identifier, getExpr.Expr)
	var calleeValue EnvironmentData
	switch getExpr.Expr.NodeName() {
	case "Identifier":
		val, err := e.identifier(getExpr.Expr)
		if err != nil {
			return nil, err
		}
		calleeValue = val

	case "FunctionCall":
		val, err := e.functionCall(getExpr.Expr)

		if err != nil {
			return nil, err
		}

		calleeValue = val
	// TODO FOR Literal Data
	default:
		return nil, e.NewError(e.NewTokenFromLine(node.GetLine()), JoError.DefaultError, "unknown callee")
	}

	if calleeValue == nil {
		return nil, e.NewError(identifier.Token, JoError.ReferenceError, fmt.Sprintf("can't access property `%s` from a null data", identifier.Value))
	}

	switch calleeValue.Type() {
	case Struct:
		struct_ := calleeValue.(*StructData)
		v, err := struct_.env.GetOne(identifier.Value)
		if err != nil {
			return nil, e.NewError(identifier.Token, JoError.DefaultError, fmt.Sprintf("method/attribute `%s` not defined", identifier.Value))
		}
		return v, nil
	case Function:
		// id := getExpr.Expr.(*parser.Identifier)

		// fun.FunctionDecl.Identifier
		return nil, e.NewError(identifier.Token, JoError.DefaultError, fmt.Sprintf("can't access property `%s` from a function declaration", identifier.Value))
	case StructDecl:
		return nil, e.NewError(identifier.Token, JoError.DefaultError, fmt.Sprintf("can't access property `%s` from a struct declaration", identifier.Value))
	// TODO FOR Literal Data
	default:
		return nil, e.NewError(e.NewTokenFromLine(node.GetLine()), JoError.DefaultError, "unknown callee")
	}
}
