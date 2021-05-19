package interpreter

import (
	"fmt"
	"math"

	"github.com/polyscone/knight/value"
)

// Exp will raise the LHS value to the power of the RHS value.
// The LHS value must be an integer, and the RHS will be converted to an
// integer if it isn't one already.
//
// If LHS is zero then RHS must be positive.
func (i *Interpreter) Exp(lhs, rhs value.Value) (value.Value, error) {
	if lhs, ok := lhs.(*value.Int); ok {
		rhs := rhs.AsInt().Value
		if lhs.Value == 0 && rhs < 0 {
			return nil, fmt.Errorf("cannot raise %v to a negative power", lhs)
		}

		return value.NewInt(int(math.Pow(float64(lhs.Value), float64(rhs)))), nil
	}

	return nil, fmt.Errorf("cannot raise %s to %s", lhs, rhs)
}
