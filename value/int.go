package value

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/polyscone/knight/build"
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
	return NewString(strconv.Itoa(i.Value))
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
	return strconv.Itoa(i.Value)
}

// NewInt will return a runtime Bool value that wraps the given int.
func NewInt(i int) *Int {
	if !build.Reckless {
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
