package interpreter

import (
	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

// Or returns either the LHS value if it was truthy, or evaluates and returns
// the RHS expression if the LHS value was falsey.
// This means that when the LHS value is true it short-circuits.
func (i *Interpreter) Or(lhs value.Value, rhs ast.Node) (value.Value, error) {
	if lhs.AsBool().Value {
		return lhs, nil
	}

	return i.eval(rhs)
}
