package interpreter

import "github.com/polyscone/knight/value"

// Call will evaluate the given value as an expression.
// It only expects to be given a block value, but will evaluate any other
// expression as undefined behaviour as well.
func (i *Interpreter) Call(arg value.Value) (value.Value, error) {
	return i.eval(arg.AsExpr())
}
