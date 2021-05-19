package interpreter

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/build"
	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

var programs = struct {
	sync.Mutex
	data map[*value.String]ast.Program
}{data: make(map[*value.String]ast.Program)}

// Eval will execute the program in the given value.
func (i *Interpreter) Eval(val value.Value) (value.Value, error) {
	s := val.AsString()

	if !build.Reckless {
		programs.Lock()
	}

	if program, ok := programs.data[s]; ok {
		if !build.Reckless {
			programs.Unlock()
		}

		return i.eval(program.Expression)
	}

	l := lexer.New()
	p := parser.New(l, i.globals)
	r := strings.NewReader(s.Value)

	program, err := p.Parse(r)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	programs.data[s] = program

	if !build.Reckless {
		programs.Unlock()
	}

	return i.eval(program.Expression)
}
