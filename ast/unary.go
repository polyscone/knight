package ast

import (
	"fmt"

	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

// Unary represents an AST node that describes a unary operation.
type Unary struct {
	value.Expr

	Op   token.Kind
	Node Node
}

// Dump prints a string form of Unary for testing.
// The Knight spec doesn't actually require AST nodes like this to print anything
// but this implementation does it anyway.
func (u Unary) Dump() string {
	return fmt.Sprintf("Unary(%s, %v)", u.Op, u.Node)
}

// String prints a string form of Unary as an s-expression for testing.
func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Op, u.Node)
}

// NewUnary returns a Unary AST node.
func NewUnary(op token.Kind, node Node) Node {
	return &Unary{
		Op:   op,
		Node: node,
	}
}
