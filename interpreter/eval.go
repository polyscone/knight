package interpreter

import (
	"strings"
	"sync"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/options"
	"github.com/polyscone/knight/value"
)

var programs = struct {
	sync.Mutex
	data map[*value.String]ast.Program
}{data: make(map[*value.String]ast.Program)}

// Eval will execute the program in the given value.
func (i *Interpreter) Eval(val value.Value) (value.Value, error) {
	s := val.AsString()

	if !options.Reckless {
		programs.Lock()
	}

	if program, ok := programs.data[s]; ok {
		if !options.Reckless {
			programs.Unlock()
		}

		return i.eval(program.Root)
	}

	program, err := i.parser.Parse(i.globals, strings.NewReader(s.Value))
	if err != nil {
		return nil, err
	}

	programs.data[s] = program

	if !options.Reckless {
		programs.Unlock()
	}

	return i.eval(program.Root)
}
