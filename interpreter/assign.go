package interpreter

import "github.com/polyscone/knight/value"

// Assign checks that the LHS expression is a Global AST node, and if it is
// evaluates the the RHS expression and assigns the resulting value to it.
func (i *Interpreter) Assign(lhs *value.Global, rhs value.Value) (value.Value, error) {
	lhs.Value = rhs

	return lhs.Value, nil
}
