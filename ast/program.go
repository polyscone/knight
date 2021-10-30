package ast

import (
	"github.com/polyscone/knight/value"
)

// Program represents a valid program as an AST.
type Program struct {
	Globals *value.GlobalStore
	Root    Node
}

// String prints a string form of Program as an s-expression for testing.
func (p Program) String() string {
	return p.ASTString("sexp")
}

// ASTString returns a string representation of the AST in the requested style.
func (p Program) ASTString(style string) string {
	return SprintNode(style, "program", p.Root.ASTString(style))
}
