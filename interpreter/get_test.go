package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		str   string
		start int
		count int
		want  string
	}{
		{"start and count zero", "foo", 0, 0, ""},
		{"start at the end and count zero", "foo", 3, 0, ""},
		{"start at the beginning and count length of string", "foo", 0, 3, "foo"},
		{"get middle letter", "foo", 1, 1, "o"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			str := value.NewString(tc.str)
			start := value.NewInt(tc.start)
			count := value.NewInt(tc.count)
			result, err := interpreter.New(nil).Get(str, start, count)
			if err != nil {
				t.Fatal(err)
			}

			if want := value.NewString(tc.want); !value.Equal(result, want) {
				t.Errorf("want %v (%p), got %v (%p)", want, want, result, result)
			}
		})
	}
}
