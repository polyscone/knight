package value_test

import (
	"testing"

	"github.com/polyscone/knight/value"
)

func TestAtoi(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		str  string
		want int
	}{
		{"zero", "0", 0},
		{"zero with multiple digits", "00000", 0},
		{"negative zero", "-0", 0},
		{"empty string", "", 0},
		{"minus sign only", "-", 0},
		{"identifier only", "a", 0},
		{"positive integer", "123", 123},
		{"negative integer", "-123", -123},
		{"positive integer with leading spaces", "       456", 456},
		{"negative integer with leading spaces", "       -456", -456},
		{"positive integer with trailing spaces", "789      ", 789},
		{"negative integer with trailing spaces", "-789      ", -789},
		{"positive integer with leading and trailing spaces", "    14785239      ", 14785239},
		{"negative integer with leading and trailing spaces", "    -14785239      ", -14785239},
		{"integer interrupted by identifier", "78a9", 78},
		{"negative integer interrupted by a second minus sign", "-78-9", -78},
		{"positive integer interrupted by a minus sign", "78-9", 78},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := value.Atoi(tc.str); got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
