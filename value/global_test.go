package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestGlobalDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		identifier string
		value      value.Value
		want       string
	}{
		{"x integer", "x", value.NewInt(1), `Global("x", Number(1))`},
		{"y string", "y", value.NewString("foo"), `Global("y", String(foo))`},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			globals := value.NewGlobalStore()
			global := globals.New(tc.identifier)

			global.Value = tc.value

			if got := global.Dump(); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
