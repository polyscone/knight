package interpreter

import (
	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

// Block wraps the given expression into a new block value and returns it.
// The wrapped expression is not evaluated until the block is passed as an
// argument to CALL.
func (i *Interpreter) Block(expr ast.Node) (value.Value, error) {
	return value.NewBlock(expr), nil
}
