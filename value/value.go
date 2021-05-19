package value

// Value represents a runtime value that can be converted to other types.
type Value interface {
	Expression

	AsBool() *Bool
	AsInt() *Int
	AsString() *String
	AsExpr() Expression
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
		if !ok {
			return false
		}

		if lhs.tag == 0 || rhs.tag == 0 {
			return lhs.Value == rhs.Value
		}

		return lhs.tag == rhs.tag
	}

	return lhs == rhs
}
