package value

import "fmt"

// Expression represents any runtime Knight expression.
// Runtime values, which implement the Value interface, also implement this
// interface as they are a slightly more concrete form of an expression.
type Expression interface {
	fmt.Stringer

	IsExpr() bool
}

// Value represents a runtime value that can be converted to other types.
type Value interface {
	Expression

	AsBool() *Bool
	AsInt() *Int
	AsString() *String
	AsExpr() Expression
	Dump() string
}

type Expr struct{}

// IsExpr indicates that the implemting type can be used as a
// runtime Knight expression.
func (e Expr) IsExpr() bool {
	return true
}

// Equal checks to see if two value are equal to each other.
//
// Equality is determined by both the type and the value.
// If LHS and RHS are of different types then they will never be equal.
func Equal(lhs, rhs Value) bool {
	switch lhs := lhs.(type) {
	case *Int:
		rhs, ok := rhs.(*Int)

		return ok && lhs.Value == rhs.Value
	case *String:
		rhs, ok := rhs.(*String)

		return ok && (lhs == rhs || lhs.Value == rhs.Value)
	}

	return lhs == rhs
}
