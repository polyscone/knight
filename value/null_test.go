package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestNullConversions(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		value      value.Value
		wantBool   value.Value
		wantInt    value.Value
		wantString value.Value
	}{
		{"null", value.NewNull(), value.NewBool(false), value.NewInt(0), value.NewString("null")},
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

func TestNullDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr value.Expression
		want string
	}{
		{"null", value.NewNull(), "Null()"},
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
