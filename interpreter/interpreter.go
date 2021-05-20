package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

var stdin = bufio.NewReader(os.Stdin)

// Parser should build a valid AST from some source code.
type Parser interface {
	Parse(globals *value.GlobalStore, r io.ByteScanner) (ast.Program, error)
}

// Interpreter is an implementation of a tree-walk interpreter than can execute
// any valid Knight AST.
type Interpreter struct {
	globals *value.GlobalStore
	parser  Parser
}

// Execute will walk the given program's AST executing nodes as it goes.
func (i *Interpreter) Execute(program ast.Program) (value.Value, error) {
	// Pre-allocate the entire range of internable ints and their string conversions
	for i := value.MinInternInt; i <= value.MaxInternInt; i++ {
		_ = value.NewInt(i)
	}

	return i.eval(program.Expression)
}

func (i *Interpreter) eval(expr value.Expression) (value.Value, error) {
	switch v := expr.(type) {
	case *value.Block:
		return v, nil
	case *value.Bool:
		return v, nil
	case *value.Int:
		return v, nil
	case *value.String:
		return v, nil
	case *value.Null:
		return v, nil
	case *value.Global:
		if v.Value == nil {
			return nil, fmt.Errorf("attempted to access undefined variable %v", v)
		}

		if _, ok := v.Value.(*value.Block); ok {
			return i.eval(v.Value)
		}

		return v.Value, nil
	case *ast.Call:
		switch v.Letter {
		case 'B':
			return i.Block(v.Args[0])
		case 'C':
			arg, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Call(arg)
		case 'D':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Dump(os.Stdout, val)
		case 'E':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Eval(val)
		case 'G':
			str, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			start, err := i.eval(v.Args[1])
			if err != nil {
				return nil, err
			}

			count, err := i.eval(v.Args[2])
			if err != nil {
				return nil, err
			}

			return i.Get(str, start, count)
		case 'I':
			condition, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.If(condition, v.Args[1], v.Args[2])
		case 'L':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Length(val)
		case 'O':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Output(os.Stdout, val)
		case 'P':
			return i.Prompt(stdin)
		case 'Q':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.Quit(val)
		case 'R':
			return i.Random()
		case 'S':
			str, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			start, err := i.eval(v.Args[1])
			if err != nil {
				return nil, err
			}

			count, err := i.eval(v.Args[2])
			if err != nil {
				return nil, err
			}

			replacement, err := i.eval(v.Args[3])
			if err != nil {
				return nil, err
			}

			return i.Substitute(str, start, count, replacement)
		case 'W':
			return i.While(v.Args[0], v.Args[1])
		}

		if v.Name == "XD" {
			fmt.Print(v.Args[0])

			return value.NewNull(), nil
		}

		return nil, fmt.Errorf("undefined function: %v", v.Name)
	case *ast.Unary:
		val, err := i.eval(v.Value)
		if err != nil {
			return nil, err
		}

		switch v.Op {
		case token.Not:
			return i.Not(val)
		case token.Noop:
			return i.Noop(val)
		case token.System:
			return i.System(val)
		default:
			return nil, fmt.Errorf("unknown unary operator: %s", v)
		}
	case *ast.Binary:
		var lhs value.Value
		var rhs value.Value

		// Assignment shouldn't evaluate its LHS because in that case we want
		// to set the value, not get it
		//
		// Chain conditionally evaluates its RHS, so we do it in the case
		// rather than here to simplify the error checks
		if v.Op != token.Assign && v.Op != token.Chain {
			var err error
			if lhs, err = i.eval(v.LHS); err != nil {
				return nil, err
			}

			// The logical and/or operations can short-circuit, so we need to
			// evaluate the RHS conditionally
			if v.Op != token.And && v.Op != token.Or {
				if rhs, err = i.eval(v.RHS); err != nil {
					return nil, err
				}
			}
		}

		switch v.Op {
		case token.Add:
			return i.Add(lhs, rhs)
		case token.Sub:
			return i.Sub(lhs, rhs)
		case token.Mul:
			return i.Mul(lhs, rhs)
		case token.Div:
			return i.Div(lhs, rhs)
		case token.Mod:
			return i.Mod(lhs, rhs)
		case token.Exp:
			return i.Exp(lhs, rhs)
		case token.Less:
			return i.Less(lhs, rhs)
		case token.Greater:
			return i.Greater(lhs, rhs)
		case token.And:
			return i.And(lhs, v.RHS)
		case token.Or:
			return i.Or(lhs, v.RHS)
		case token.Equal:
			return i.Equal(lhs, rhs)
		case token.Chain:
			if _, err := i.eval(v.LHS); err != nil {
				return nil, err
			}

			return i.eval(v.RHS)
		case token.Assign:
			global, ok := v.LHS.(*value.Global)
			if !ok {
				return nil, fmt.Errorf("cannot assign to %s", v.LHS)
			}

			rhs, err := i.eval(v.RHS)
			if err != nil {
				return nil, err
			}

			global.Value = rhs

			return rhs, nil
		default:
			return nil, fmt.Errorf("unknown binary operator: %s", v)
		}
	default:
		return nil, fmt.Errorf("unknown expression: %s", expr)
	}
}

// New returns an initialised Interpreter that can be used to execute programs
// that are represented as an AST.
func New(globals *value.GlobalStore, parser Parser) *Interpreter {
	return &Interpreter{
		globals: globals,
		parser:  parser,
	}
}
