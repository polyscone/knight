package interpreter_test

import (
	"testing"

	"github.com/polyscone/knight/interpreter"
	"github.com/polyscone/knight/value"
)

func TestSubstitute(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		str     string
		start   int
		count   int
		replace string
		want    string
	}{
		{"substitute nothing with nothing", "foo", 0, 0, "", "foo"},
		{"substitute nothing with bar", "foo", 0, 0, "bar", "barfoo"},
		{"substitute f with bar", "foo", 0, 1, "bar", "baroo"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			str := value.NewString(tc.str)
			start := value.NewInt(tc.start)
			count := value.NewInt(tc.count)
			replace := value.NewString(tc.replace)
			result, err := interpreter.New(nil, nil).Substitute(str, start, count, replace)
			if err != nil {
				t.Fatal(err)
			}

			if want := value.NewString(tc.want); !value.Equal(result, want) {
				t.Errorf("want %v (%p), got %v (%p)", want, want, result, result)
			}
		})
	}
}
