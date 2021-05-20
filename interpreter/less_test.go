package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestLess(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  value.Value
		want bool
	}{
		{"false is less than true", value.NewBool(false), value.NewBool(true), true},
		{"true is not less than false", value.NewBool(true), value.NewBool(false), false},
		{"1 is less than 2", value.NewInt(1), value.NewInt(2), true},
		{"2 is not less than 1", value.NewInt(2), value.NewInt(1), false},
		{"a is less than b", value.NewString("a"), value.NewString("b"), true},
		{"b is not less than a", value.NewString("b"), value.NewString("a"), false},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Less(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if want := value.NewBool(tc.want); !value.Equal(result, want) {
				t.Errorf("want %v (%p), got %v (%p)", want, want, result, result)
			}
		})
	}
}
