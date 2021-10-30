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
	return u.ASTString("sexp")
}

// ASTString returns a string representation of the AST in the requested style.
func (u Unary) ASTString(style string) string {
	return SprintNode(style, u.Op.String(), u.Node.ASTString(style))
}

// NewUnary returns a Unary AST node.
func NewUnary(op token.Kind, node Node) Node {
	return &Unary{
		Op:   op,
		Node: node,
	}
}
