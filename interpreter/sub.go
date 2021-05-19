package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Sub returns the result of the subtraction of the RHS value from the LHS value.
// The LHS value must be an integer, and the RHS value will be converted to an
// integer if it isn't one already.
func (i *Interpreter) Sub(lhs, rhs value.Value) (value.Value, error) {
	if lhs, ok := lhs.(*value.Int); ok {
		if rhs, ok := rhs.(*value.Int); ok && rhs.Value == 0 {
			return lhs, nil
		}

		return value.NewInt(lhs.Value - rhs.AsInt().Value), nil
	}

	return nil, fmt.Errorf("cannot subtract %s from %s", rhs, lhs)
}
