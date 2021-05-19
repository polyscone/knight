package interpreter

import (
	"os"

	"github.com/polyscone/knight/value"
)

// Quit will exit the program with the integer conversion of the given value.
func (i *Interpreter) Quit(val value.Value) (value.Value, error) {
	os.Exit(val.AsInt().Value)

	return nil, nil
}
