package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestNoop(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		value value.Value
		want  value.Value
	}{
		{"noop bool", value.NewBool(true), value.NewBool(true)},
		{"noop int", value.NewInt(123), value.NewInt(123)},
		{"noop string", value.NewString("foo"), value.NewString("foo")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil).Noop(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
