package interpreter

import "github.com/polyscone/knight/value"

// Noop does nothing and just returns the value it was given.
func (i *Interpreter) Noop(val value.Value) (value.Value, error) {
	return val, nil
}
