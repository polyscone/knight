package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// ASCII returns either the string representation of an integer or the integer
// representation of the first character of a string.
func (i *Interpreter) ASCII(val value.Value) (value.Value, error) {
	switch val.(type) {
	case *value.Int:
		return value.NewString(string(rune(val.AsInt().Value))), nil
	case *value.String:
		return value.NewInt(int(val.AsString().Value[0])), nil
	}

	return nil, fmt.Errorf("unknown ascii argument %v", val)
}
