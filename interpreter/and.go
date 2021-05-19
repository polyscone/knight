package interpreter

import "github.com/polyscone/knight/value"

// And returns either the LHS value if it was falsey, or evaluates and returns
// the RHS expression if the LHS value was truthy.
// This means that when the LHS value is false it short-circuits.
func (i *Interpreter) And(lhs value.Value, rhsExpr value.Expression) (value.Value, error) {
	if !lhs.AsBool().Value {
		return lhs, nil
	}

	return i.eval(rhsExpr)
}
