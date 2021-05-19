package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Mod returns the remainder of the division of the LHS value by the RHS value.
// The LHS value must be an integer, and the RHS value will be converted to an
// integer if it isn't one already.
//
// The RHS value must be a positive number.
func (i *Interpreter) Mod(lhs, rhs value.Value) (value.Value, error) {
	if lhs, ok := lhs.(*value.Int); ok {
		rhs := rhs.AsInt().Value
		if rhs <= 0 {
			return nil, fmt.Errorf("cannot modulo by %v", rhs)
		}

		return value.NewInt(lhs.Value % rhs), nil
	}

	return nil, fmt.Errorf("cannot divide %s by %s", lhs, rhs)
}
