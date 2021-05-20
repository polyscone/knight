package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

func TestCall(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr value.Expression
		want value.Value
	}{
		{"call wrapped int", value.NewInt(1), value.NewInt(1)},
		{"call wrapped addition", ast.NewBinary(token.Add, value.NewInt(1), value.NewInt(2)), value.NewInt(3)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			i := interpreter.New(nil, nil)
			block, err := i.Block(tc.expr)
			if err != nil {
				t.Fatal(err)
			}

			result, err := i.Call(block)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
