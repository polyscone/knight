package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Assign checks that the LHS expression is a Global AST node, and if it is
// evaluates the the RHS expression and assigns the resulting value to it.
func (i *Interpreter) Assign(lhsExpr, rhsExpr value.Expression) (value.Value, error) {
	global, ok := lhsExpr.(*value.Global)
	if !ok {
		return nil, fmt.Errorf("cannot assign to %s", lhsExpr)
	}

	rhs, err := i.eval(rhsExpr)
	if err != nil {
		return nil, err
	}

	global.Value = rhs

	return rhs, nil
}
