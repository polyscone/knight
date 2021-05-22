package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/polyscone/knight/ast"
	"github.com/polyscone/knight/token"
	"github.com/polyscone/knight/value"
)

var builtinArities = map[byte]int{
	'B': 1,
	'C': 1,
	'D': 1,
	'E': 1,
	'G': 3,
	'I': 3,
	'L': 1,
	'O': 1,
	'P': 0,
	'Q': 1,
	'R': 0,
	'S': 4,
	'W': 2,
}

// Lexer should provide a stream of tokens from some source code.
type Lexer interface {
	Load(r io.ByteScanner)
	Peek() token.Token
	Consume() (token.Token, error)
}

// Parser holds the state for a parser than can transform a stream of Knight
// tokens into an AST.
type Parser struct {
	lexer   Lexer
	globals *value.GlobalStore
}

// Parse will load the source code int the given byte scanner into its lexer and
// build an AST from the resulting token stream.
func (p *Parser) Parse(globals *value.GlobalStore, r io.ByteScanner) (ast.Program, error) {
	p.lexer.Load(r)
	p.globals = globals

	program := ast.Program{Globals: p.globals}

	expr, err := p.parseExpr()
	if err != nil {
		return program, err
	}

	program.Root = expr

	return program, nil
}

func (p *Parser) parseExpr() (ast.Node, error) {
	tok, err := p.lexer.Consume()
	if err != nil {
		return ast.Invalid, err
	}

	switch tok.Kind {
	case token.Integer:
		i, err := strconv.Atoi(tok.Lexeme)
		if err != nil {
			return ast.Invalid, err
		}

		return ast.NewInt(i), nil
	case token.String:
		return ast.NewString(tok.Lexeme), nil
	case token.True, token.False:
		return ast.NewBool(tok.Kind == token.True), nil
	case token.Null:
		return ast.NewNull(), nil
	case token.Not, token.Noop, token.System:
		value, err := p.parseExpr()
		if err != nil {
			return ast.Invalid, err
		}

		return ast.NewUnary(tok.Kind, value), nil
	case token.And,
		token.Or,
		token.Add,
		token.Sub,
		token.Mul,
		token.Div,
		token.Mod,
		token.Less,
		token.Greater,
		token.Assign,
		token.Equal,
		token.Exp,
		token.Chain:

		lhs, err := p.parseExpr()
		if err != nil {
			return ast.Invalid, err
		}

		rhs, err := p.parseExpr()
		if err != nil {
			return ast.Invalid, err
		}

		return ast.NewBinary(tok.Kind, lhs, rhs), nil
	case token.Variable:
		return ast.NewGlobal(p.globals.New(tok.Lexeme)), nil
	case token.Call:
		letter := tok.Lexeme[0]
		arity, ok := builtinArities[letter]
		if !ok {
			if tok.Lexeme == "XD" {
				arity = 1
			} else {
				return ast.Invalid, fmt.Errorf("unexpected function %q", tok.Lexeme)
			}
		}

		args := make([]ast.Node, arity)
		for i := 0; i < arity; i++ {
			arg, err := p.parseExpr()
			if err != nil {
				return ast.Invalid, err
			}

			args[i] = arg
		}

		return ast.NewCall(tok.Lexeme, args), nil
	default:
		return ast.Invalid, fmt.Errorf("unexpected token: %s", tok)
	}
}

// New returns a new initialised Parser.
func New(lexer Lexer) *Parser {
	return &Parser{lexer: lexer}
}
