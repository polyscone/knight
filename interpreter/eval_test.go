package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

func TestEval(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		source string
		want   value.Value
	}{
		{"bool true", "T", value.NewBool(true)},
		{"bool false", "F", value.NewBool(false)},
		{"integer addition", "+ 1 2", value.NewInt(3)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			p := parser.New(l)
			result, err := interpreter.New(nil, p).Eval(value.NewString(tc.source))
			if err != nil {
				t.Fatal(err)
			}

			if result != tc.want {
				t.Errorf("want %v, got %v", tc.want, result)
			}
		})
	}
}
