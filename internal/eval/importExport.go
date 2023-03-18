package eval

import (
	"errors"
	"path/filepath"
	"strings"

	Node "github.com/overlorddamygod/jo/pkg/parser/node"
)

func (e *Evaluator) Import(node Node.Node) (EnvironmentData, error) {
	import_, ok := node.(*Node.ImportStatement)

	if !ok {
		return nil, errors.New("node is not an import statement")
	}

	path, err := filepath.Abs(import_.File.Literal + ".jo")
	if err != nil {
		panic(err)
	}

	evaluator, err := Import(path)

	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(path)
	key := strings.Split(fileName, ".")[0]

	if evaluator.exportData == nil {
		return nil, errors.New("no export data found in file " + fileName)
	}

	if import_.Alias != nil {
		key = import_.Alias.Literal
	}

	err = e.global.DefineOne(key, evaluator.exportData)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Evaluator) Export(node Node.Node) (EnvironmentData, error) {
	export, ok := node.(*Node.ExportStatement)

	if !ok {
		return nil, errors.New("node is not an export statement")
	}

	exportData, err := e.EvalExpression(export.Expr)
	if err != nil {
		return nil, err
	}

	e.exportData = exportData

	return nil, nil
}
