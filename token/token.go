package token

import (
	"fmt"
	"strconv"
)

// Tokens kinds.
const (
	EOF Kind = iota
	Err
	Unknown

	Call
	Variable
	Integer
	String
	True
	False
	Null

	And
	Or
	Not
	Add
	Sub
	Mul
	Div
	Mod
	Less
	Greater
	Assign
	Equal
	Exp
	System
	Chain
	Noop
)

// Kind is used to describe the kind of token.
type Kind byte

// String returns a string representation of the kind.
func (k Kind) String() string {
	switch k {
	case EOF:
		return "EOF"
	case Err:
		return "error"
	case Unknown:
		return "unknown"
	case Call:
		return "call"
	case Variable:
		return "variable"
	case Integer:
		return "integer"
	case String:
		return "string"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"
	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"
	case Add:
		return "add"
	case Sub:
		return "sub"
	case Mul:
		return "mul"
	case Div:
		return "div"
	case Mod:
		return "mod"
	case Less:
		return "less"
	case Greater:
		return "greater"
	case Assign:
		return "assign"
	case Equal:
		return "equal"
	case Exp:
		return "exp"
	case System:
		return "system"
	case Chain:
		return "chain"
	case Noop:
		return "noop"
	}

	return strconv.Itoa(int(k))
}

// Token represents an atomic piece of code.
type Token struct {
	Kind   Kind
	Lexeme string
	Line   int
	Col    int
}

// String returns a string representation of the Token along with line and
// column information.
func (t Token) String() string {
	return fmt.Sprintf("%v:%v %#q (%s)", t.Line, t.Col, t.Lexeme, t.Kind)
}
