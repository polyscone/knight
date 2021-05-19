package interpreter

import "github.com/polyscone/knight/value"

// Not returns the negation of the boolean conversion of the given value.
func (i *Interpreter) Not(val value.Value) (value.Value, error) {
	return value.NewBool(!val.AsBool().Value), nil
}
