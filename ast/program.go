package ast

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Program represents a valid program as an AST.
type Program struct {
	Globals *value.GlobalStore
	Root    Node
}

// String prints a string form of Program as an s-expression for testing.
func (p Program) String() string {
	return fmt.Sprintf("(program %s)", p.Root)
}
