package interpreter

import (
	"fmt"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

// Call will evaluate the given value as an expression.
// It only expects to be given a block value, but will evaluate any other
// expression as undefined behaviour as well.
func (i *Interpreter) Call(arg value.Value) (value.Value, error) {
	switch v := arg.AsExpr().(type) {
	case *value.Bool:
		return v, nil
	case *value.Int:
		return v, nil
	case *value.String:
		return v, nil
	case *value.Null:
		return v, nil
	case ast.Node:
		return i.eval(v)
	}

	return nil, fmt.Errorf("unknown call argument %v", arg)
}
