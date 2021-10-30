package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestChain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		lhs  ast.Node
		rhs  ast.Node
	}{
		{"lhs int, rhs int", value.NewInt(1), value.NewInt(2)},
		{"lhs string, rhs string", value.NewString("foo"), value.NewString("bar")},
		{"lhs int, rhs string", value.NewInt(1), value.NewString("bar")},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).Chain(tc.lhs, tc.rhs)
			if err != nil {
				t.Fatal(err)
			}

			// To compare interfaces one side must be assignable to the other
			// Since Value and Node aren't assignable to each other we convert
			// them both to the empty interface
			// This makes them assignable to each other and allows for comparison
			// of their dynamic values
			if interface{}(result) != interface{}(tc.rhs) {
				t.Errorf("want %v (%p), got %v (%p)", tc.rhs, tc.rhs, result, result)
			}
		})
	}
}
