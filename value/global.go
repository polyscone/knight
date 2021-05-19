package value

import (
	"fmt"
	"sync"

	"github.com/polyscone/knight/build"
)

// Global represents a global name/value pair.
type Global struct {
	Name  string
	Value Value
}

// Dump prints a string form of Global for testing.
// A global doesn't actually need to dump anything according to the Knight spec
// but this implementation prints a representation of it anyway.
func (g Global) Dump() string {
	return fmt.Sprintf("Global(%q, %v)", g.Name, g.Value.Dump())
}

// String prints a string form of Global as an s-expression for testing.
func (g Global) String() string {
	return fmt.Sprintf("(var %q)", g.Name)
}

// GlobalStore holds the state for a group of global name/value pairs.
// Globals from one store should not be mixed with another.
type GlobalStore struct {
	sync.Mutex
	data map[string]*Global
}

// New returns a Global with the given name.
// If the global doesn't exist yet then it is created.
//
// Globals are always cached by name in the given store, so calling this
// function multiple times with the same store and name will return the
// same global object.
func (gs *GlobalStore) New(name string) *Global {
	if !build.Reckless {
		gs.Lock()
		defer gs.Unlock()
	}

	if v, ok := gs.data[name]; ok {
		return v
	}

	v := Global{Name: name}

	gs.data[name] = &v

	return &v
}

func NewGlobalStore() *GlobalStore {
	return &GlobalStore{data: make(map[string]*Global)}
}
