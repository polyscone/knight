package interpreter

import (
	"fmt"
	"io"

	"github.com/polyscone/knight/value"
)

// Output prints the string conversion of the given value.
// If the given value ends with a backslash (\) then the backslash is omitted
// and no newline is printed, otherwise a newline is printed with the value.
func (i *Interpreter) Output(w io.Writer, val value.Value) (value.Value, error) {
	out := val.AsString().Value

	if out != "" && out[len(out)-1] == '\\' {
		fmt.Fprint(w, out[:len(out)-1])
	} else {
		fmt.Fprintln(w, out)
	}

	return value.NewNull(), nil
}
