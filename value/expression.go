package value

import "fmt"

// Expression represents any runtime Knight expression.
// Runtime values, which implement the Value interface, also implement this
// interface as they are a slightly more concrete form of an expression.
type Expression interface {
	fmt.Stringer

	Dump() string
}
