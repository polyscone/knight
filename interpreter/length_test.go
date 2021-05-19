package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestLength(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		str  string
		want value.Value
	}{
		{"empty 0", "", value.NewInt(0)},
		{"foo 3", "foo", value.NewInt(3)},
		{"foo bar 7", "foo bar", value.NewInt(7)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil).Length(value.NewString(tc.str))
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
