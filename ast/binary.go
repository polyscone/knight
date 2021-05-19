package ast

import (
	"fmt"

	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

// Binary represents an AST node that expects to hole two
// non-nil value.Expressions.
type Binary struct {
	Op  token.Kind
	LHS value.Expression
	RHS value.Expression
}

// Dump prints a string form of Binary for testing.
// The Knight spec doesn't actually require AST nodes like this to print anything
// but this implementation does it anyway.
func (b Binary) Dump() string {
	return fmt.Sprintf("Binary(%s, %v, %v)", b.Op, b.LHS.Dump(), b.RHS.Dump())
}

// String prints a string form of Binary as an s-expression for testing.
func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Op, b.LHS, b.RHS)
}

// NewBinary returns a binary AST node that expects to hold two
// non-nil value.Expressions.
func NewBinary(op token.Kind, lhs, rhs value.Expression) *Binary {
	return &Binary{
		Op:  op,
		LHS: lhs,
		RHS: rhs,
	}
}
