package eval

import (
	Node "github.com/overlorddamygod/jo/pkg/parser/node"
	"github.com/overlorddamygod/jo/pkg/stdio"
)

func Print(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	output := ""
	for i, arg := range arguments {
		exp, err := e.EvalExpression(arg)

		if err != nil {
			return nil, err
		}

		if i > 0 {
			output += " "
		}

		output += exp.GetString()
	}
	stdio.Io.Println(output)
	return nil, nil
}

func Input(e *Evaluator, name string, arguments []Node.Node) (EnvironmentData, error) {
	arg1 := arguments[0]
	arg, err := e.EvalExpression(arg1)

	if err != nil {
		return nil, err
	}
	argLiteral := arg.(LiteralData)

	stdio.Io.Print(argLiteral.GetString())

	text := stdio.Io.Input()
	return StringLiteral(text), nil
}
