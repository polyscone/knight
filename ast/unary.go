package ast

import (
	"fmt"

	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

// Unary represents an AST node that describes a unary operation.
type Unary struct {
	Op    token.Kind
	Value value.Expression
}

// Dump prints a string form of Unary for testing.
// The Knight spec doesn't actually require AST nodes like this to print anything
// but this implementation does it anyway.
func (u Unary) Dump() string {
	return fmt.Sprintf("Unary(%s, %v)", u.Op, u.Value.Dump())
}

// String prints a string form of Unary as an s-expression for testing.
func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Op, u.Value)
}

// NewUnary returns a Unary AST node.
func NewUnary(op token.Kind, value value.Expression) *Unary {
	return &Unary{
		Op:    op,
		Value: value,
	}
}
