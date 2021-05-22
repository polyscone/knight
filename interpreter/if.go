package interpreter

import (
	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

// If will evaluate and return the consequence value if the boolean result of
// the condition value is true, otherwise it will evaluate and return the
// alternative value.
func (i *Interpreter) If(condition value.Value, consequence, alternative ast.Node) (value.Value, error) {
	if condition.AsBool().Value {
		return i.eval(consequence)
	}

	return i.eval(alternative)
}
