package value

import "github.com/polyscone/knight/ast"

var null = &Null{}

// Int represents a runtime null value.
type Null struct{}

// AsBool converts the caller to a runtime Bool representation of Null.
func (n *Null) AsBool() *Bool {
	return _false
}

// AsInt converts the caller to a runtime Int representation of Null.
func (n *Null) AsInt() *Int {
	return zero
}

// AsString converts the caller to a runtime String representation of Null.
func (n *Null) AsString() *String {
	return nullString
}

// AsExpr returns the value itself as an Expression interface implementation.
func (n *Null) AsExpr() Expression {
	return n
}

// Dump prints a string form of Null for testing.
func (n *Null) Dump() string {
	return "Null()"
}

// String prints a string form of the Null as an s-expression for testing.
// The AsString method should be used to convert a value to a runtime String.
func (n *Null) String() string {
	return n.ASTString(ast.StyleSexpr)
}

// ASTString returns a string representation of the AST in the requested style.
func (n *Null) ASTString(style ast.Style) string {
	return "null"
}

// NewNull will return a runtime Null value.
func NewNull() *Null {
	return null
}
