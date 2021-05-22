package value_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/value"
)

func TestIntConversions(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		value      value.Value
		wantBool   value.Value
		wantInt    value.Value
		wantString value.Value
	}{
		{"zero", value.NewInt(0), value.NewBool(false), value.NewInt(0), value.NewString("0")},
		{"one", value.NewInt(1), value.NewBool(true), value.NewInt(1), value.NewString("1")},
		{"negative", value.NewInt(-123), value.NewBool(true), value.NewInt(-123), value.NewString("-123")},
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

func TestIntDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr ast.Node
		want string
	}{
		{"zero", value.NewInt(0), "Number(0)"},
		{"one", value.NewInt(1), "Number(1)"},
		{"negative", value.NewInt(-123), "Number(-123)"},
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
