package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  value.Value
		want value.Value
	}{
		{"lhs positive int, rhs positive int", value.NewInt(1), value.NewInt(2), value.NewInt(3)},
		{"lhs negative int, rhs negative int", value.NewInt(-1), value.NewInt(-2), value.NewInt(-3)},
		{"lhs positive int, rhs string", value.NewInt(25), value.NewString("10"), value.NewInt(35)},
		{"lhs string, rhs positive int", value.NewString("foo"), value.NewInt(3), value.NewString("foo3")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Add(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
