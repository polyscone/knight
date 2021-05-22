package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestIf(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		condition   value.Value
		consequence ast.Node
		alternative ast.Node
		want        value.Value
	}{
		{"condition truthy eval consequence", value.NewBool(true), value.NewInt(1), value.NewInt(2), value.NewInt(1)},
		{"condition falsey eval alternative", value.NewBool(false), value.NewInt(1), value.NewInt(2), value.NewInt(2)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).If(tc.condition, tc.consequence, tc.alternative)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
