package ast

import "github.com/polyscone/knight/token"

// Binary represents an AST node that expects to hold two non-nil ast.Nodes.
type Binary struct {
	Op  token.Kind
	LHS Node
	RHS Node
}

// String prints a string form of Binary as an s-expression for testing.
func (b Binary) String() string {
	return b.ASTString(StyleSexpr)
}

// ASTString returns a string representation of the AST in the requested style.
func (b Binary) ASTString(style Style) string {
	return SprintNode(
		style,
		b.Op.String(),
		b.LHS.ASTString(style),
		b.RHS.ASTString(style),
	)
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
