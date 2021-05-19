package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestDiv(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  value.Value
		want value.Value
	}{
		{"lhs positive int, rhs positive int", value.NewInt(10), value.NewInt(2), value.NewInt(5)},
		{"lhs positive int, rhs positive int, truncate", value.NewInt(9), value.NewInt(2), value.NewInt(4)},
		{"lhs negative int, rhs negative int", value.NewInt(-10), value.NewInt(-5), value.NewInt(2)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil).Div(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
