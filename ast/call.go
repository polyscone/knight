package ast

import (
	"fmt"
	"strings"

	"github.com/polyscone/knight/value"
)

// Call represents an AST node that describes a function call.
// The Args are represented as a four-element array, but since not all function
// calls will use all four elements of the array the NArgs field should be
// used to see how many arguments were actually provided.
type Call struct {
	Name   string
	Letter byte
	Args   [4]value.Expression
	NArgs  int
}

// Dump prints a string form of Call for testing.
// The Knight spec doesn't actually require AST nodes like this to print anything
// but this implementation does it anyway.
func (c Call) Dump() string {
	args := make([]string, 0, c.NArgs)
	for i := 0; i < c.NArgs; i++ {
		args = append(args, c.Args[i].Dump())
	}

	return fmt.Sprintf("Call(%q, %c, Args(%s))", c.Name, c.Letter, strings.Join(args, ", "))
}

// String prints a string form of Call as an s-expression for testing.
func (c Call) String() string {
	if c.NArgs == 0 {
		return c.Name
	}

	args := make([]string, 0, c.NArgs)
	for i := 0; i < c.NArgs; i++ {
		args = append(args, c.Args[i].String())
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
func NewCall(name string, args []value.Expression) *Call {
	if len(args) > 4 {
		panic("too many args")
	}

	var argsArr [4]value.Expression
	for i, arg := range args {
		argsArr[i] = arg
	}

	return &Call{
		Name:   name,
		Letter: name[0],
		Args:   argsArr,
		NArgs:  len(args),
	}
}
