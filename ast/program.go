package ast

// Program represents a valid program as an AST.
type Program struct {
	Root Node
}

// String prints a string form of Program as an s-expression for testing.
func (p Program) String() string {
	return p.ASTString(StyleSexpr)
}

// ASTString returns a string representation of the AST in the requested style.
func (p Program) ASTString(style Style) string {
	return SprintNode(style, "program", p.Root.ASTString(style))
}
