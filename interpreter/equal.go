package interpreter

import (
	"github.com/polyscone/knight/value"
)

// Equal returns the boolean result of an equality comparison between the
// LHS and RHS values.
//
// For LHS and RHS to be considered equal they must both be of the same
// value and type.
func (i *Interpreter) Equal(lhs, rhs value.Value) (value.Value, error) {
	return value.NewBool(value.Equal(lhs, rhs)), nil
}
