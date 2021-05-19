package interpreter

import "github.com/polyscone/knight/value"

// Length returns the length in bytes of the given value when converted to a string.
func (i *Interpreter) Length(val value.Value) (value.Value, error) {
	return value.NewInt(len(val.AsString().Value)), nil
}
