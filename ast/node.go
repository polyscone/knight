package ast

import "github.com/polyscone/knight/value"

// Invalid represents an invalid AST node that, if part of a larger AST, should
// render the entire tree invalid.
var Invalid Node

// Node is a wrapper around the lower level AST node.
// It also wraps concrete values in the case of leaf nodes.
type Node interface {
	value.Expression
}
