package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestNot(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		value value.Value
		want  value.Value
	}{
		{"true becomes false", value.NewBool(true), value.NewBool(false)},
		{"false becomes true", value.NewBool(false), value.NewBool(true)},
		{"1 becomes false", value.NewInt(1), value.NewBool(false)},
		{"0 becomes true", value.NewInt(0), value.NewBool(true)},
		{"string foo becomes false", value.NewString("foo"), value.NewBool(false)},
		{"empty string becomes true", value.NewString(""), value.NewBool(true)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Not(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
