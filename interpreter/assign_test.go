package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestAssign(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		global string
		rhs    value.Value
		want   value.Value
	}{
		{"assign string to x", "x", value.NewString("foo"), value.NewString("foo")},
		{"assign int to x", "x", value.NewInt(2), value.NewInt(2)},
		{"assign int to y", "y", value.NewInt(1), value.NewInt(1)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			g := value.NewGlobalStore()
			global := g.New(tc.global)

			_, err := interpreter.New(g, nil).Assign(global, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(global.Value, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, global.Value, global.Value)
			}
		})
	}
}
