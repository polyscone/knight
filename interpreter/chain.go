package interpreter

import (
	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

// Chain evaluates both the given LHS and RHS expressions, but if the LHS
// expressions returns an error execution will stop.
// The return value is always the result of the RHS expression.
func (i *Interpreter) Chain(lhsExpr, rhsExpr ast.Node) (value.Value, error) {
	if _, err := i.eval(lhsExpr); err != nil {
		return nil, err
	}

	return i.eval(rhsExpr)
}
