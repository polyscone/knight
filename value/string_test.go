package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestStringConversions(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		value      value.Value
		wantBool   value.Value
		wantInt    value.Value
		wantString value.Value
	}{
		{"empty", value.NewString(""), value.NewBool(false), value.NewInt(0), value.NewString("")},
		{"single word", value.NewString("foo"), value.NewBool(true), value.NewInt(0), value.NewString("foo")},
		{"space separated", value.NewString("foo bar"), value.NewBool(true), value.NewInt(0), value.NewString("foo bar")},
		{"starts with a number", value.NewString("123 foo bar"), value.NewBool(true), value.NewInt(123), value.NewString("123 foo bar")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.value.AsBool(); !value.Equal(got, tc.wantBool) {
				t.Errorf("want %v (%p), got %v (%p)", tc.wantBool, tc.wantBool, got, got)
			}

			if got := tc.value.AsInt(); !value.Equal(got, tc.wantInt) {
				t.Errorf("want %v (%p), got %v (%p)", tc.wantInt, tc.wantInt, got, got)
			}

			if got := tc.value.AsString(); !value.Equal(got, tc.wantString) {
				t.Errorf("want %v (%p), got %v (%p)", tc.wantString, tc.wantString, got, got)
			}
		})
	}
}

func TestStringDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr value.Expression
		want string
	}{
		{"empty", value.NewString(""), "String()"},
		{"single word", value.NewString("foo"), "String(foo)"},
		{"space separated", value.NewString("foo bar"), "String(foo bar)"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.expr.Dump(); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
