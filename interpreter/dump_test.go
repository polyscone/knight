package interpreter_test

import (
	"bytes"
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		value value.Value
		want  string
	}{
		{"bool true", value.NewBool(true), "Bool(true)"},
		{"bool false", value.NewBool(false), "Bool(false)"},
		{"int 123", value.NewInt(123), "Number(123)"},
		{`string "foo"`, value.NewString("foo"), "String(foo)"},
		{`string "foo bar"`, value.NewString("foo bar"), "String(foo bar)"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			_, err := interpreter.New(nil, nil).Dump(&buf, tc.value)
			if err != nil {
				t.Fatal(err)
			}

			if got := buf.String(); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
