package interpreter

import (
	"fmt"
	"io"

	"github.com/polyscone/knight/value"
)

// Dump will print the debug/test string output for a given value.
func (i *Interpreter) Dump(w io.Writer, val value.Value) (value.Value, error) {
	fmt.Fprint(w, val.Dump())

	return value.NewNull(), nil
}
