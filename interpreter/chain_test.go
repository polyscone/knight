package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestChain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Expression
		rhs  value.Expression
	}{
		{"lhs int, rhs int", value.NewInt(1), value.NewInt(2)},
		{"lhs string, rhs string", value.NewString("foo"), value.NewString("bar")},
		{"lhs int, rhs string", value.NewInt(1), value.NewString("bar")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil).Chain(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if result != tc.rhs {
				t.Errorf("want %v (%p), got %v (%p)", tc.rhs, tc.rhs, result, result)
			}
		})
	}
}