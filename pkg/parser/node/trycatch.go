package node

import (
	"fmt"

	"github.com/overlorddamygod/jo/pkg/lexer"
)

type TryCatchStatement struct {
	name     string
	Try      *Block
	CatchVar *Identifier
	Catch    *Block
}

func NewTryCatchStatement(try *Block, catchVar *Identifier, catch *Block) *TryCatchStatement {
	return &TryCatchStatement{
		name:     TRY_CATCH,
		Try:      try,
		CatchVar: catchVar,
		Catch:    catch,
		// ElseBlock: Else,
	}
}

func (a *TryCatchStatement) NodeName() string {
	return a.name
}

func (a TryCatchStatement) GetLine() int {
	return a.Try.GetLine()
}

func (a *TryCatchStatement) Print() {
	fmt.Println(a.name)
	// a.Print()
}

type ThrowStatement struct {
	name       string
	token      *lexer.Token
	Expression Node
}

func NewThrowStatement(token *lexer.Token, expression Node) *ThrowStatement {
	return &ThrowStatement{
		name:       THROW,
		token:      token,
		Expression: expression,
	}
}

func (a *ThrowStatement) NodeName() string {
	return a.name
}

func (a *ThrowStatement) Print() {
	fmt.Println(a.name)
	a.Expression.Print()
	// a.Expression.Print()
}

func (a ThrowStatement) GetLine() int {
	return a.token.GetLine()
}
