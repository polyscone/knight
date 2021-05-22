package ast

import (
	"fmt"

	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

// Binary represents an AST node that expects to hole two
// non-nil ast.Nodes.
type Binary struct {
	value.Expr

	Op  token.Kind
	LHS Node
	RHS Node
}

// String prints a string form of Binary as an s-expression for testing.
func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Op, b.LHS, b.RHS)
}

// NewBinary returns a binary AST node that expects to hold two
// non-nil ast.Nodes.
func NewBinary(op token.Kind, lhs, rhs Node) Node {
	return &Binary{
		Op:  op,
		LHS: lhs,
		RHS: rhs,
	}
}
