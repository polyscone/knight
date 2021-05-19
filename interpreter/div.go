package interpreter

import (
	"errors"
	"fmt"

	"github.com/polyscone/knight/value"
)

// Div returns the result of dividing the LHS value by the RHS value.
// LHS must be an integer, and the RHS value will be converted to an integer
// if it isn't one already.
//
// Dividing by zero will return an error.
func (i *Interpreter) Div(lhs, rhs value.Value) (value.Value, error) {
	if lhs, ok := lhs.(*value.Int); ok {
		rhs := rhs.AsInt().Value
		if rhs == 0 {
			return nil, errors.New("cannot divide by 0")
		}

		return value.NewInt(lhs.Value / rhs), nil
	}

	return nil, fmt.Errorf("cannot divide %s by %s", lhs, rhs)
}
