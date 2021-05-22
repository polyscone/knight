package value

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/polyscone/knight/build"
)

// MaxInternStringLength defines the maximum length of an interned string.
// Any strings that go over this length will not be interned.
const MaxInternStringLength = 32

var strs = struct {
	sync.Mutex
	data map[string]*String
}{data: make(map[string]*String)}

var concats = struct {
	sync.Mutex
	data map[int]*String
}{data: make(map[int]*String)}

var substrs = struct {
	sync.Mutex
	data map[int]*String
}{data: make(map[int]*String)}

var substrsSlow = struct {
	sync.Mutex
	data map[int]map[int]*String
}{data: make(map[int]map[int]*String)}

var intStrs = struct {
	sync.Mutex
	data map[int]*String
}{data: make(map[int]*String)}

var (
	nullString  = NewString("null")
	trueString  = NewString("true")
	falseString = NewString("false")
)

// String represents a runtime string value.
type String struct {
	Expr

	Value string
	tag   int
}

// AsBool converts the caller to a false Bool runtime value if the caller's
// value is an empty string, or to a true Bool value otherwise.
func (s *String) AsBool() *Bool {
	if s.Value != "" {
		return _true
	}

	return _false
}

// AsInt converts the caller to a runtime Int value.
func (s *String) AsInt() *Int {
	return NewInt(Atoi(s.Value))
}

// AsString returns the caller without modification.
func (s *String) AsString() *String {
	return s
}

// AsExpr returns the value itself as an Expression interface implementation.
func (s *String) AsExpr() Expression {
	return s
}

// Dump prints a string form of String for testing.
func (s *String) Dump() string {
	return fmt.Sprintf("String(%v)", s.Value)
}

// String prints a string form of the String as an s-expression for testing.
// The AsString method should be used to convert a value to a runtime String.
func (s *String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

// NewString will return a runtime String value that wraps the given string.
func NewString(s string) *String {
	if len(s) > MaxInternStringLength {
		return NewUniqueString(s)
	}

	if !build.Reckless {
		strs.Lock()
		defer strs.Unlock()
	}

	if v, ok := strs.data[s]; ok {
		return v
	}

	v := String{
		Value: s,
		tag:   len(strs.data) + 1,
	}

	strs.data[s] = &v

	return &v
}

// NewUniqueString will return a runtime String value that wraps the given
// string, but the returned object will always be newly allocated and never interned.
func NewUniqueString(s string) *String {
	return &String{Value: s}
}

// NewSubString will return a runtime String value that wraps the given string
// from the given start and end indices.
func NewSubString(s *String, start, end int) *String {
	if s.tag == 0 {
		return NewString(s.Value[start:end])
	}

	if start <= math.MaxUint16 && end <= math.MaxUint16 {
		if !build.Reckless {
			substrs.Lock()
			defer substrs.Unlock()
		}

		key := s.tag<<32 | start<<16 | end

		if v, ok := substrs.data[key]; ok {
			return v
		}

		v := NewString(s.Value[start:end])

		substrs.data[key] = v

		return v
	}

	if !build.Reckless {
		substrsSlow.Lock()
		defer substrsSlow.Unlock()
	}

	key1 := s.tag
	key2 := start<<32 | end

	if _, ok := substrsSlow.data[key1]; !ok {
		substrsSlow.data[key1] = make(map[int]*String)
	}

	if v, ok := substrsSlow.data[key1][key2]; ok {
		return v
	}

	v := NewString(s.Value[start:end])

	substrsSlow.data[key1][key2] = v

	return v
}

// NewConcatString will return a runtime String value that is the concatenation
// of both given runtime String values.
func NewConcatString(lhs, rhs *String) *String {
	if lhs.tag == 0 || rhs.tag == 0 {
		return NewString(lhs.Value + rhs.Value)
	}

	if !build.Reckless {
		concats.Lock()
		defer concats.Unlock()
	}

	key := lhs.tag<<32 | rhs.tag

	if v, ok := concats.data[key]; ok {
		return v
	}

	v := NewString(lhs.Value + rhs.Value)

	concats.data[key] = v

	return v
}

// NewIntString will return a runtime String value that is the string
// representation of the given integer.
func NewIntString(i int) *String {
	if !build.Reckless {
		intStrs.Lock()
		defer intStrs.Unlock()
	}

	if v, ok := intStrs.data[i]; ok {
		return v
	}

	v := NewString(strconv.Itoa(i))

	intStrs.data[i] = v

	return v
}
