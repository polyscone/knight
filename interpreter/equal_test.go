package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  value.Value
		want bool
	}{
		{"true equals true", value.NewBool(true), value.NewBool(true), true},
		{"true does not equal false", value.NewBool(true), value.NewBool(false), false},
		{"1 equals 1", value.NewInt(1), value.NewInt(1), true},
		{"1 does not equal 2", value.NewInt(1), value.NewInt(2), false},
		{"foo equals foo", value.NewString("foo"), value.NewString("foo"), true},
		{"foo does not equal bar", value.NewString("foo"), value.NewString("bar"), false},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Equal(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if want := value.NewBool(tc.want); !value.Equal(result, want) {
				t.Errorf("want %v (%p), got %v (%p)", want, want, result, result)
			}
		})
	}
}
