package value

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/polyscone/knight/options"
)

// These values define the range of ints to intern.
// Any ints that fall outside this range will not be interned.
const (
	MinInternInt = -10
	MaxInternInt = math.MaxUint8
)

// These values are created immediately so they can be used by other values, like
// Bool and Null, without forcing them to call NewInt.
var (
	zero = NewInt(0)
	one  = NewInt(1)
)

var ints = struct {
	sync.Mutex
	data map[int]*Int
}{data: make(map[int]*Int)}

// Int represents a runtime integer value.
type Int struct {
	Expr

	Value int
}

// AsBool converts the caller to a false Bool runtime value if the caller's
// value is 0, or to a true Bool value otherwise.
func (i *Int) AsBool() *Bool {
	return NewBool(i.Value != 0)
}

// AsInt returns the caller without modification.
func (i *Int) AsInt() *Int {
	return i
}

// AsString converts the caller to a runtime String representation of its value.
func (i *Int) AsString() *String {
	return NewIntString(i.Value)
}

// AsExpr returns the value itself as an Expression interface implementation.
func (i *Int) AsExpr() Expression {
	return i
}

// Dump prints a string form of Int for testing.
func (i *Int) Dump() string {
	return fmt.Sprintf("Number(%v)", i.Value)
}

// String prints a string form of the Int as an s-expression for testing.
// The AsString method should be used to convert a value to a runtime String.
func (i *Int) String() string {
	return i.ASTString("sexp")
}

// ASTString returns a string representation of the AST in the requested style.
func (i *Int) ASTString(style string) string {
	return strconv.Itoa(i.Value)
}

// NewInt will return a runtime Int value that wraps the given int.
func NewInt(i int) *Int {
	if i < MinInternInt || i > MaxInternInt {
		return NewUniqueInt(i)
	}

	if !options.Reckless {
		ints.Lock()
		defer ints.Unlock()
	}

	if v, ok := ints.data[i]; ok {
		return v
	}

	v := Int{Value: i}

	ints.data[i] = &v

	return &v
}

// NewUniqueInt will return a runtime Int value that wraps the given int, but the
// returned object will always be newly allocated and never interned.
func NewUniqueInt(i int) *Int {
	return &Int{Value: i}
}
