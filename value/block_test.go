package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestBlockConversions(t *testing.T) {
	t.Parallel()

	t.Run("as bool panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("wanted panic")
			}
		}()

		value.NewBlock(value.NewInt(0)).AsBool()
	})

	t.Run("as int panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("wanted panic")
			}
		}()

		value.NewBlock(value.NewInt(0)).AsInt()
	})

	t.Run("as string panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("wanted panic")
			}
		}()

		value.NewBlock(value.NewInt(0)).AsString()
	})

	t.Run("as expr", func(t *testing.T) {
		t.Parallel()

		expr := value.NewInt(0)

		if value.NewBlock(expr).AsExpr() != expr {
			t.Error("want block expression to equal original expression")
		}
	})
}

func TestBlockDump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		expr value.Expression
		want string
	}{
		{"bool true", value.NewBool(true), "Block(Bool(true))"},
		{"bool false", value.NewBool(false), "Block(Bool(false))"},
		{"int 0", value.NewInt(0), "Block(Number(0))"},
		{"int 1", value.NewInt(1), "Block(Number(1))"},
		{`string "foo"`, value.NewString("foo"), "Block(String(foo))"},
		{`string "bar"`, value.NewString("bar"), "Block(String(bar))"},
		{"null", value.NewNull(), "Block(Null())"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := value.NewBlock(tc.expr).Dump(); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
