package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestAnd(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  value.Value
		rhs  ast.Node
		want value.Value
	}{
		{"lhs true, rhs int", value.NewBool(true), value.NewInt(1), value.NewInt(1)},
		{"lhs false, rhs int", value.NewBool(false), value.NewInt(1), value.NewBool(false)},
		{"lhs truthy int, rhs int", value.NewInt(1), value.NewInt(2), value.NewInt(2)},
		{"lhs falsey int, rhs int", value.NewInt(0), value.NewInt(2), value.NewInt(0)},
		{"lhs truthy string, rhs int", value.NewString("foo"), value.NewInt(2), value.NewInt(2)},
		{"lhs falsey string, rhs int", value.NewString(""), value.NewInt(2), value.NewString("")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).And(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
