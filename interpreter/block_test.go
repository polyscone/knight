package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

func TestBlock(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr ast.Node
	}{
		{"wrap int", value.NewInt(1)},
		{"wrap addition", ast.NewBinary(token.Add, value.NewInt(1), value.NewInt(2))},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			block, err := interpreter.New(nil, nil).Block(tc.expr)
			if err != nil {
				t.Fatal(err)
			}

			if got := block.AsExpr(); got != tc.expr {
				t.Errorf("want %v (%p), got %v (%p)", tc.expr, tc.expr, got, got)
			}
		})
	}
}
