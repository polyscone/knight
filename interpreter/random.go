package interpreter

import (
	"math/rand"

	"github.com/polyscone/knight/value"
)

// Random returns a pseudo-random integer value.
func (i *Interpreter) Random() (value.Value, error) {
	//nolint:gosec // there's no need for a CSPRNG here
	return value.NewInt(int(rand.Int63())), nil
}
