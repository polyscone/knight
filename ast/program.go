package ast

import (
	"fmt"

	"github.com/polyscone/knight/value"
)

// Program represents the root of an AST.
type Program struct {
	Globals    *value.GlobalStore
	Expression value.Expression
}

// String prints a string form of Program as an s-expression for testing.
func (p Program) String() string {
	return fmt.Sprintf("(program %s)", p.Expression)
}
