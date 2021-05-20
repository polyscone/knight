package interpreter_test

import (
	"bytes"
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestOutput(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		value value.Value
		want  string
	}{
		{"bool true", value.NewBool(true), "true\n"},
		{"bool false", value.NewBool(false), "false\n"},
		{"integer 123", value.NewInt(123), "123\n"},
		{"string foo bar", value.NewString("foo bar"), "foo bar\n"},
		{"string foo bar escape no newline", value.NewString("foo bar\\"), "foo bar"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			_, err := interpreter.New(nil, nil).Output(&buf, tc.value)
			if err != nil {
				t.Fatal(err)
			}

			if got := buf.String(); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
