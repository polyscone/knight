package ast

import (
	"fmt"
	"strings"

	"github.com/polyscone/knight/value"
)

// Call represents an AST node that describes a function call.
type Call struct {
	value.Expr

	Name   string
	Letter byte
	Args   []Node
}

// Dump prints a string form of Call for testing.
// The Knight spec doesn't actually require AST nodes like this to print anything
// but this implementation does it anyway.
func (c Call) Dump() string {
	args := make([]string, len(c.Args))
	for i := range c.Args {
		args[i] = c.Args[i].String()
	}

	return fmt.Sprintf("Call(%q, %c, Args(%s))", c.Name, c.Letter, strings.Join(args, ", "))
}

// String prints a string form of Call as an s-expression for testing.
func (c Call) String() string {
	if len(c.Args) == 0 {
		return c.Name
	}

	args := make([]string, len(c.Args))
	for i := range c.Args {
		args[i] = c.Args[i].String()
	}

	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(c.Name)
	sb.WriteString(" ")
	sb.WriteString(strings.Join(args, " "))
	sb.WriteString(")")

	return sb.String()
}

// NewCall returns a Call AST node that represents a call with the given
// arguments to a function.
// Args are provided as a slice, but there is an upper-limit of four args.
func NewCall(name string, args []Node) Node {
	if len(args) > 4 {
		panic("too many args")
	}

	return &Call{
		Name:   name,
		Letter: name[0],
		Args:   args,
	}
}
