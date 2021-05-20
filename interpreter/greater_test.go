package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestGreater(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  value.Value
		want bool
	}{
		{"true is greater than false", value.NewBool(true), value.NewBool(false), true},
		{"false is not greater than true", value.NewBool(false), value.NewBool(true), false},
		{"2 is greater than 1", value.NewInt(2), value.NewInt(1), true},
		{"1 is not greater than 2", value.NewInt(1), value.NewInt(2), false},
		{"b is greater than a", value.NewString("b"), value.NewString("a"), true},
		{"a is not greater than b", value.NewString("a"), value.NewString("b"), false},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Greater(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if want := value.NewBool(tc.want); !value.Equal(result, want) {
				t.Errorf("want %v (%p), got %v (%p)", want, want, result, result)
			}
		})
	}
}
