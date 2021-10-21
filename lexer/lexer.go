package lexer

import (
	"errors"
	"io"
	"strings"

	"github.com/polyscone/knight/token"
)

const eof byte = 0

// ErrUnknownToken is returned with every consumed token that isn't specified
// by the language in some way.
var ErrUnknownToken = errors.New("unknown token")

type predicate func(byte) bool

type position struct {
	line int
	col  int
}

// Lexer holds the state for a scanner that tokenises Knight source code.
type Lexer struct {
	r       io.ByteScanner
	lastPos position
	pos     position
	curr    token.Token
	next    token.Token
	err     error
}

// Load resets the state of the lexer and prepares it to tokenise the source
// code in the given byte scanner.
func (l *Lexer) Load(r io.ByteScanner) {
	l.r = r
	l.pos = position{line: 1, col: 0}
	l.curr = token.Token{Kind: token.EOF}
	l.next = token.Token{Kind: token.EOF}
	l.err = nil

	if l.r != nil {
		_, _ = l.Consume()
	}
}

// Peek returns the next token in the stream without consuming anything.
func (l *Lexer) Peek() token.Token {
	return l.next
}

// Consume returns the next token in the stream along with any error that might
// have occurred during the scan.
func (l *Lexer) Consume() (token.Token, error) {
	curr, err := l.next, l.err

	l.curr = curr
	l.next, l.err = l.consume()

	return curr, err
}

func (l *Lexer) consume() (token.Token, error) {
start:
	pos := l.pos
	if _, err := l.readWhile(isWhitespace); err != nil {
		return l.newError(err, pos)
	}

	r, err := l.peek()
	if err != nil {
		return l.newError(err, l.pos)
	}

	if isComment(r) {
		pos := l.pos
		if _, err := l.readWhile(notNewline); err != nil {
			return l.newError(err, pos)
		}

		goto start
	}

	pos = l.pos
	switch {
	case isEOF(r):
		return l.newToken(token.EOF, "", pos)
	case isPunc(r):
		lexeme, err := l.read()
		if err != nil {
			return l.newError(err, pos)
		}

		switch r {
		case '&':
			return l.newToken(token.And, "&", pos)
		case '|':
			return l.newToken(token.Or, "|", pos)
		case '!':
			return l.newToken(token.Not, "!", pos)
		case '+':
			return l.newToken(token.Add, "+", pos)
		case '-':
			return l.newToken(token.Sub, "-", pos)
		case '*':
			return l.newToken(token.Mul, "*", pos)
		case '/':
			return l.newToken(token.Div, "/", pos)
		case '%':
			return l.newToken(token.Mod, "%", pos)
		case '<':
			return l.newToken(token.Less, "<", pos)
		case '>':
			return l.newToken(token.Greater, ">", pos)
		case '=':
			return l.newToken(token.Assign, "=", pos)
		case '?':
			return l.newToken(token.Equal, "?", pos)
		case '^':
			return l.newToken(token.Exp, "^", pos)
		case '`':
			return l.newToken(token.System, "`", pos)
		case ';':
			return l.newToken(token.Chain, ";", pos)
		case ':':
			return l.newToken(token.Noop, ":", pos)
		}

		return l.newUnknown(string(lexeme), pos)
	case isWordFuncStart(r):
		lexeme, err := l.readWhile(isWordFunc)
		if err != nil {
			return l.newError(err, pos)
		}

		switch r {
		case 'T':
			return l.newToken(token.True, lexeme, pos)
		case 'F':
			return l.newToken(token.False, lexeme, pos)
		case 'N':
			return l.newToken(token.Null, lexeme, pos)
		}

		return l.newToken(token.Call, lexeme, pos)
	case isIdentStart(r):
		lexeme, err := l.readWhile(isIdent)
		if err != nil {
			return l.newError(err, pos)
		}

		return l.newToken(token.Variable, lexeme, pos)
	case isDigit(r):
		lexeme, err := l.readWhile(isDigit)
		if err != nil {
			return l.newError(err, pos)
		}

		return l.newToken(token.Integer, lexeme, pos)
	case isString(r):
		// Discard open quote
		if _, err := l.read(); err != nil {
			return l.newError(err, pos)
		}

		lexeme, err := l.readUntil(r)
		if err != nil {
			return l.newError(err, pos)
		}

		// Discard close quote
		if _, err := l.read(); err != nil {
			return l.newError(err, pos)
		}

		return l.newToken(token.String, lexeme, pos)
	}

	lexeme, err := l.readWhile(notWhitespace)
	if err != nil {
		return l.newError(err, pos)
	}

	return l.newUnknown(lexeme, pos)
}

func (l *Lexer) newToken(kind token.Kind, lexeme string, pos position) (token.Token, error) {
	tok := token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Line:   pos.line,
		Col:    pos.col,
	}

	return tok, nil
}

func (l *Lexer) newUnknown(lexeme string, pos position) (token.Token, error) {
	tok, _ := l.newToken(token.Unknown, lexeme, pos)

	return tok, ErrUnknownToken
}

func (l *Lexer) newError(err error, pos position) (token.Token, error) {
	tok, _ := l.newToken(token.Err, "", pos)

	return tok, err
}

func (l *Lexer) peek() (byte, error) {
	r, err := l.r.ReadByte()
	//nolint:errorlint // io.EOF is never wrapped
	if err == io.EOF {
		return eof, nil
	}
	if err != nil {
		return r, err
	}

	return r, l.r.UnreadByte()
}

func (l *Lexer) read() (byte, error) {
	r, err := l.r.ReadByte()
	//nolint:errorlint // io.EOF is never wrapped
	if err == io.EOF {
		return eof, nil
	}
	if err != nil {
		return r, err
	}

	l.lastPos = l.pos
	if r == '\n' {
		l.pos.line++
		l.pos.col = 1
	} else {
		l.pos.col++
	}

	return r, nil
}

func (l *Lexer) unread() error {
	if err := l.r.UnreadByte(); err != nil {
		return err
	}

	l.pos, l.lastPos = l.lastPos, l.pos

	return nil
}

func (l *Lexer) readWhile(f predicate) (string, error) {
	var sb strings.Builder

	for {
		r, err := l.read()
		if err != nil {
			return sb.String(), err
		}
		if r == eof {
			return sb.String(), nil
		}

		if !f(r) {
			return sb.String(), l.unread()
		}

		sb.WriteByte(r)
	}
}

func (l *Lexer) readUntil(r byte) (string, error) {
	return l.readWhile(func(current byte) bool { return current != r })
}

// New returns a new initialised Lexer.
func New() *Lexer {
	l := Lexer{}

	l.Load(nil)

	return &l
}

func isNewline(r byte) bool {
	return r == '\n'
}

func notNewline(r byte) bool {
	return !isNewline(r)
}

func isWhitespace(r byte) bool {
	switch r {
	case ' ', '\t', '\n', '\r', ':', '(', ')', '{', '}', '[', ']':
		return true
	}

	return false
}

func notWhitespace(r byte) bool {
	return !isWhitespace(r)
}

func isEOF(r byte) bool {
	return r == eof
}

func isPunc(r byte) bool {
	switch r {
	case '!', '%', '&', '*', '+', '-', '/', ':', ';', '<', '=', '>', '?', '^', '`', '|':
		return true
	}

	return false
}

func isWordFuncStart(r byte) bool {
	switch r {
	case 'B', 'C', 'D', 'E', 'F', 'G', 'I', 'L', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'W', 'X':
		return true
	}

	return false
}

func isWordFunc(r byte) bool {
	return (r >= 'A' && r <= 'Z') || r == '_'
}

func isIdentStart(r byte) bool {
	return (r >= 'a' && r <= 'z') || r == '_'
}

func isIdent(r byte) bool {
	return isIdentStart(r) || isDigit(r)
}

func isDigit(r byte) bool {
	return r >= '0' && r <= '9'
}

func isString(r byte) bool {
	return r == '"' || r == '\''
}

func isComment(r byte) bool {
	return r == '#'
}
