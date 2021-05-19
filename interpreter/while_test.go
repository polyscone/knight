package interpreter_test

import (
	"strings"
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

func TestWhile(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		source string
		want   value.Value
	}{
		{"count to 10", `; (= i 0) ; WHILE (< i 10) (= i (+ i 1)) : i`, value.NewInt(10)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			g := value.NewGlobalStore()
			p := parser.New(l, g)
			program, err := p.Parse(strings.NewReader(tc.source))
			if err != nil {
				t.Fatal(err)
			}

			result, err := interpreter.New(g).Execute(program)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
