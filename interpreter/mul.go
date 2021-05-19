package interpreter

import (
	"fmt"
	"strings"

	"github.com/polyscone/knight/value"
)

// Mul returns the product of the LHS and RHS values.
//
// If the LHS value is an integer then the result will be the multiplication of
// the LHS value with the integer conversion of the RHS value.
//
// If the LHS value is a string then the result will be the LSH value repeated
// the number of times specified by the integer conversion of the RHS value.
func (i *Interpreter) Mul(lhs, rhs value.Value) (value.Value, error) {
	switch lhs := lhs.(type) {
	case *value.Int:
		return value.NewInt(lhs.Value * rhs.AsInt().Value), nil
	case *value.String:
		count := rhs.AsInt().Value
		if count < 0 {
			return nil, fmt.Errorf("invalid string repeat count %v", count)
		}

		return value.NewString(strings.Repeat(lhs.Value, count)), nil
	}

	return nil, fmt.Errorf("cannot subtract %s from %s", rhs, lhs)
}
