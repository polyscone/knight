package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestASCII(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		value value.Value
		want  value.Value
	}{
		{`integer 38 should convert to string "&"`, value.NewInt(38), value.NewString("&")},
		{`integer 59 should convert to string ";"`, value.NewInt(59), value.NewString(";")},
		{`integer 10 should convert to string "\n"`, value.NewInt(10), value.NewString("\n")},
		{`string "Hello" should convert to integer 72`, value.NewString("Hello"), value.NewInt(72)},
		{`string "\n" should convert to integer 10`, value.NewString("\n"), value.NewInt(10)},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := interpreter.New(nil, nil).ASCII(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			if !value.Equal(result, tc.want) {
				t.Errorf("want %v (%p), got %v (%p)", tc.want, tc.want, result, result)
			}
		})
	}
}
