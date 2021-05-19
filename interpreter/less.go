package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Less returns the boolean result of a less than (<) comparison between
// the LHS and RHS values.
//
// The LHS value can be either a boolean, integer, or string, and the RHS value
// will be converted to the same type.
func (i *Interpreter) Less(lhs, rhs value.Value) (value.Value, error) {
	switch lhs := lhs.(type) {
	case *value.Bool:
		return value.NewBool(!lhs.Value && rhs.AsBool().Value), nil
	case *value.Int:
		return value.NewBool(lhs.Value < rhs.AsInt().Value), nil
	case *value.String:
		return value.NewBool(lhs.Value < rhs.AsString().Value), nil
	}

	return nil, fmt.Errorf("cannot compare %s and %s", lhs, rhs)
}
