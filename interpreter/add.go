package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Add combines the given LHS and RHS values.
//
// If the LHS is an integer then the RHS will be converted to an integer and
// the return value will be the sum of the two.
//
// If the LHS is a string then the RHS will be converted to a string and the
// return value will be the concatenation of the two.
func (i *Interpreter) Add(lhs, rhs value.Value) (value.Value, error) {
	switch lhs := lhs.(type) {
	case *value.Int:
		if lhs.Value == 0 {
			if rhs, ok := rhs.(*value.Int); ok {
				return rhs, nil
			}

			return rhs.AsInt(), nil
		}

		rhs := rhs.AsInt()
		if rhs.Value == 0 {
			return lhs, nil
		}

		return value.NewInt(lhs.Value + rhs.AsInt().Value), nil
	case *value.String:
		if lhs.Value == "" {
			if rhs, ok := rhs.(*value.String); ok {
				return rhs, nil
			}

			return rhs.AsString(), nil
		}

		rhs := rhs.AsString()
		if rhs.Value == "" {
			return lhs, nil
		}

		return value.NewConcatString(lhs, rhs), nil
	default:
		return nil, fmt.Errorf("cannot add %s and %s", lhs, rhs)
	}
}
