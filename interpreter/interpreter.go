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
	return i.eval(program.Root)
}

func (i *Interpreter) eval(node ast.Node) (value.Value, error) {
	switch v := node.(type) {
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
			return nil, fmt.Errorf("attempted to access undefined variable %v", node)
		}

		if b, ok := v.Value.(*value.Block); ok {
			return i.eval(b.Value.(ast.Node))
		}

		return v.Value, nil
	case *ast.Call:
		switch v.Letter {
		case 'A':
			val, err := i.eval(v.Args[0])
			if err != nil {
				return nil, err
			}

			return i.ASCII(val)
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
			_ = v.Args[2]
			arg0, arg1, arg2 := v.Args[0], v.Args[1], v.Args[2]

			str, err := i.eval(arg0)
			if err != nil {
				return nil, err
			}

			start, err := i.eval(arg1)
			if err != nil {
				return nil, err
			}

			count, err := i.eval(arg2)
			if err != nil {
				return nil, err
			}

			return i.Get(str, start, count)
		case 'I':
			_ = v.Args[2]
			arg0, arg1, arg2 := v.Args[0], v.Args[1], v.Args[2]

			condition, err := i.eval(arg0)
			if err != nil {
				return nil, err
			}

			return i.If(condition, arg1, arg2)
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
			_ = v.Args[3]
			arg0, arg1, arg2, arg3 := v.Args[0], v.Args[1], v.Args[2], v.Args[3]

			str, err := i.eval(arg0)
			if err != nil {
				return nil, err
			}

			start, err := i.eval(arg1)
			if err != nil {
				return nil, err
			}

			count, err := i.eval(arg2)
			if err != nil {
				return nil, err
			}

			replacement, err := i.eval(arg3)
			if err != nil {
				return nil, err
			}

			return i.Substitute(str, start, count, replacement)
		case 'W':
			_ = v.Args[1]
			arg0, arg1 := v.Args[0], v.Args[1]

			return i.While(arg0, arg1)
		}

		return nil, fmt.Errorf("undefined function: %v", v.Name)
	case *ast.Unary:
		val, err := i.eval(v.Node)
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
		var lhs, rhs value.Value
		var err error

		if v.Op != token.Assign {
			if lhs, err = i.eval(v.LHS); err != nil {
				return nil, err
			}
		}

		// The logical and/or operations can short-circuit, so we need to
		// evaluate the RHS conditionally
		if v.Op != token.And && v.Op != token.Or {
			if rhs, err = i.eval(v.RHS); err != nil {
				return nil, err
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
			return rhs, nil
		case token.Assign:
			global, ok := v.LHS.(*value.Global)
			if !ok {
				return nil, fmt.Errorf("cannot assign to %s", v.LHS)
			}

			return i.Assign(global, rhs)
		default:
			return nil, fmt.Errorf("unknown binary operator: %s", v)
		}
	default:
		return nil, fmt.Errorf("unknown node: %s", node)
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
