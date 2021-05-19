package interpreter_test

import (
	"bytes"
	"testing"

	"github.com/polyscone/knight/interpreter"
)

func TestPrompt(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input string
		want  string
	}{
		{"read until lf and trim crlf", "foo\r\nbar\r\n", "foo"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			buf.WriteString(tc.input)

			result, err := interpreter.New(nil).Prompt(&buf)
			if err != nil {
				t.Fatal(err)
			}

			if got := result.AsString().Value; got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
