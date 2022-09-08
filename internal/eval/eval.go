package eval

import (
	"errors"
	"fmt"
	"math"

	JoError "github.com/overlorddamygod/jo/pkg/error"
	L "github.com/overlorddamygod/jo/pkg/lexer"
	"github.com/overlorddamygod/jo/pkg/parser"
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

type Evaluator struct {
	lexer       *L.Lexer
	node        []Node.Node
	global      *Environment
	environment *Environment
	current     Node.Node
	// variables   map[string]parser.LiteralValue
}

func NewEvaluator(lexer *L.Lexer, node []Node.Node) *Evaluator {
	env := NewEnvironment()
	return &Evaluator{lexer: lexer, node: node, global: env, environment: env}
}

func Init(src string) {
	_lexer := L.NewLexer(src)

	_, token, err := _lexer.Lex()
	// tokens, err := lexer.Lex()
	if err != nil {
		//stdio.Io.Println(tokens)
		stdio.Io.Error("[Lexer]\n\n" + JoError.New(_lexer, token, JoError.LexicalError, err.Error()).Error())
		return
	}
	// stdio.Io.Println(tokens)
	// return

	_parser := parser.NewParser(_lexer)

	node, err := _parser.Parse()

	// for _, s := range node {
	// 	s.Print()
	// }
	if err != nil {
		stdio.Io.Error("[Parser]\n\n" + err.Error())
		return
	}

	// // for _, s := range node {
	// // 	s.Print()
	// // }

	evaluator := NewEvaluator(_lexer, node)

	_, err = evaluator.Eval()

	if err != nil {
		stdio.Io.Error("[Evaluator]\n\n" + err.Error())
		return
	}
}

func (e *Evaluator) SetLexerNode(lexer *L.Lexer, node []Node.Node) {
	e.lexer = lexer
	e.node = node
}

func NewEvaluatorWithParent(e *Evaluator, parent *Environment) *Evaluator {
	// env := NewEnvironment()
	env := NewEnvironmentWithParent(parent)
	return &Evaluator{lexer: e.lexer, node: e.node, global: env, environment: env}
}

func (e *Evaluator) Eval() (EnvironmentData, error) {
	return e.EvalStatements(e.node)
}

func (e *Evaluator) EvalStatements(statements []Node.Node) (EnvironmentData, error) {

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

func (e *Evaluator) EvalStatement(node Node.Node) (EnvironmentData, error) {
	// e.Eval()
	// fmt.Println("___")
	// fmt.Println("NODE", node.NodeName())
	// node.Print()
	// fmt.Println("___")
	// return nil, nil
	switch node.NodeName() {
	case "VarDecl":
		varDecl := node.(*Node.VarDeclStatement)

		id := varDecl.Identifier.(*Node.Identifier)
		exp, err := e.EvalExpression(*varDecl.Expression)

		if err != nil {
			return nil, err
		}

		if _, err := e.environment.Get(id.Value); err == nil {
			return nil, e.NewError(id.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
		} else {
			e.environment.Define(id.Value, exp)
		}
	case "ASSIGNMENT":
		return e.assignment(node)
	case "FunctionCall":
		return e.functionCall(node)
	case "StructDecl":
		structD := node.(*Node.StructDeclStatement)

		id := structD.Identifier.(*Node.Identifier)

		if _, err := e.environment.Get(id.Value); err == nil {
			return nil, e.NewError(id.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", id.Value))
		} else {
			e.environment.Define(id.Value, NewStructDataDecl(*structD, e.environment))
		}
	case "IF":
		// fmt.Println("IF Start")
		ifStatement := node.(*Node.IfStatement)
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
		forStatement := node.(*Node.ForStatement)
		prev := e.current
		e.current = forStatement
		e.begin()

		data, err := e.EvalStatement(forStatement.Initial)

		if err != nil {
			return data, nil
		}

		for {
			e.begin()

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

			e.end()

			_, err = e.EvalStatement(forStatement.Expression)

			if err != nil {
				return nil, err
			}
		}
		e.end()
	case "FunctionDecl":
		functionDecl := node.(*Node.FunctionDeclStatement)

		functionName := functionDecl.Identifier.(*Node.Identifier)
		// fmt.Println("GLOBASTART-----")
		// e.global.Print()
		// fmt.Println("GLOBALEND------")
		if _, err := e.environment.Get(functionName.Value); err == nil {
			return nil, e.NewError(functionName.Token, JoError.DefaultError, fmt.Sprintf("Variable ` %s ` already defined", functionName.Value))
		} else {
			e.environment.Define(functionName.Value, NewCallableFunction(*functionDecl, e.environment, nil))
		}
		// fmt.Println("ENVSTART-----")
		// e.environment.Print()
		// fmt.Println("ENVEND------")
		// e.environment.Print()
	case "ReturnStatement":
		returnStmt := node.(*Node.ReturnStatement)

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
	case "Identifier", "BinaryExpression", "GetExpr":
		return e.EvalExpression(node)
	default:
		return nil, fmt.Errorf("unknown statement %s", node.NodeName())
	}
	return nil, nil
}

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
	return nil, nil
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
			return nil, e.NewError(e.NewTokenFromLine(getExpr.GetLine()), JoError.ReferenceError, fmt.Sprintf("Cannot assign to null data"))
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

func (e *Evaluator) functionCall(node Node.Node) (EnvironmentData, error) {
	functionCall := node.(*Node.FunctionCall)

	functionName, _ := functionCall.Identifier.(*Node.Identifier)

	var function EnvironmentData
	switch functionCall.Identifier.NodeName() {
	case "Identifier":
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
	case "GetExpr":
		_structMethod, err := e._get(functionCall.Identifier)
		if err != nil {
			return nil, err
		}

		function = _structMethod
	case "FunctionCall":
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
		if callableFunction.parent != nil {
			// TODO Struct attributes
			// fmt.Println("METHOD")
			// &callableFunction.parent.env.
			// callableFunction.
		}

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

func (e *Evaluator) NewTokenFromLine(line int) *L.Token {
	return L.NewToken(L.IDENTIFIER, "", line, 0, 0)
}

func (e *Evaluator) NewError(token *L.Token, _type JoError.JoErrorType, message string) error {
	return JoError.New(e.lexer, token, _type, message)
}
