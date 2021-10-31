package value

import (
	"fmt"

	"github.com/polyscone/knight/ast"
)

// _true and _false represent runtime values for true and false booleans.
// No other Bool values should be created at runtime.
var (
	_true  = &Bool{Value: true}
	_false = &Bool{Value: false}
)

// Bool represents a runtime boolean value.
type Bool struct {
	Value bool
}

// AsBool returns the caller without modification.
func (b *Bool) AsBool() *Bool {
	return b
}

// AsInt converts the caller to a runtime Int value.
func (b *Bool) AsInt() *Int {
	if b.Value {
		return one
	}

	return zero
}

// AsString converts the caller to a runtime String value.
func (b *Bool) AsString() *String {
	if b.Value {
		return trueString
	}

	return falseString
}

// AsExpr returns the value itself as an Expression interface implementation.
func (b *Bool) AsExpr() Expression {
	return b
}

// Dump prints a string form of Bool for testing.
func (b *Bool) Dump() string {
	return fmt.Sprintf("Bool(%v)", b.Value)
}

// String prints a string form of the Bool as an s-expression for testing.
// The AsString method should be used to convert a value to a runtime String.
func (b *Bool) String() string {
	return b.ASTString(ast.StyleSexpr)
}

// ASTString returns a string representation of the AST in the requested style.
func (b *Bool) ASTString(style ast.Style) string {
	if b.Value {
		return "true"
	}

	return "false"
}

// NewBool will return a runtime Bool value that wraps the given boolean.
func NewBool(b bool) *Bool {
	if b {
		return _true
	}

	return _false
}
