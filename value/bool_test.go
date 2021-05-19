package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestBoolConversions(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		value      value.Value
		wantBool   value.Value
		wantInt    value.Value
		wantString value.Value
	}{
		{"true", value.NewBool(true), value.NewBool(true), value.NewInt(1), value.NewString("true")},
		{"false", value.NewBool(false), value.NewBool(false), value.NewInt(0), value.NewString("false")},
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

func TestBoolDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr value.Expression
		want string
	}{
		{"true", value.NewBool(true), "Bool(true)"},
		{"false", value.NewBool(false), "Bool(false)"},
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
