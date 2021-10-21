package lexer_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/token"
)

func TestComments(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  string
		tokens int
	}{
		// EOF is ignored in the token count since if it doesn't exist we'll infinitely loop anyway
		{"comment only", "# This is a comment", 0},
		{"integer then comment", "1234 # This is a comment", 1},
		{"comment then newline", "# This is a comment\n", 0},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			var tokens []token.Token
			for tok, err := l.Consume(); tok.Kind != token.EOF; tok, err = l.Consume() {
				if err != nil {
					t.Fatal(err)
				}

				tokens = append(tokens, tok)
			}

			if got := len(tokens); tc.tokens != got {
				t.Errorf("want %v tokens, got %v", tc.tokens, got)
			}
		})
	}
}

func TestPuncTokens(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input string
		kind  token.Kind
	}{
		{"logical and (&)", "&", token.And},
		{"logical or (|)", "|", token.Or},
		{"negation (!)", "!", token.Not},
		{"add (+)", "+", token.Add},
		{"sub (-)", "-", token.Sub},
		{"mul (*)", "*", token.Mul},
		{"div (/)", "/", token.Div},
		{"modulo (%)", "%", token.Mod},
		{"less than (<)", "<", token.Less},
		{"greater than (>)", ">", token.Greater},
		{"assign (=)", "=", token.Assign},
		{"equal (?)", "?", token.Equal},
		{"exponentiation (^)", "^", token.Exp},
		{"system (`)", "`", token.System},
		{"chaining (;)", ";", token.Chain},
		// {"no-op (:)", ":", token.Noop},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			tok, err := l.Consume()
			if err != nil {
				t.Fatal(err)
			}

			if want := 1; tok.Line != want {
				t.Errorf("want line %v, got %v", want, tok.Line)
			}
			if tok.Kind != tc.kind {
				t.Errorf("want kind %q, got %q", tc.kind, tok.Kind)
			}
			if tok.Lexeme != tc.input {
				t.Errorf("want lexeme %q, got %q", tc.input, tok.Lexeme)
			}
		})
	}
}

func TestFunctionTokens(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input string
	}{
		{"word function ascii", "A"},
		{"word function ascii long", "ASCII"},
		{"word function block", "B"},
		{"word function block long", "BLOCK"},
		{"word function call", "C"},
		{"word function call long", "CALL"},
		{"word function dump", "D"},
		{"word function dump long", "DUMP"},
		{"word function eval", "E"},
		{"word function eval long", "EVAL"},
		{"word function get", "G"},
		{"word function get long", "GET"},
		{"word function if", "I"},
		{"word function if long", "IF"},
		{"word function length", "L"},
		{"word function length long", "LENGTH"},
		{"word function output", "O"},
		{"word function output long", "OUTPUT"},
		{"word function prompt", "P"},
		{"word function prompt long", "PROMPT"},
		{"word function quit", "Q"},
		{"word function quit long", "QUIT"},
		{"word function rand", "R"},
		{"word function rand long", "RAND"},
		{"word function substitute", "S"},
		{"word function substitute long", "SUBSTITUTE"},
		{"word function set long", "SET"},
		{"word function while", "W"},
		{"word function while long", "WHILE"},
		{"word function x", "X"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			tok, err := l.Consume()
			if err != nil {
				t.Fatal(err)
			}

			if want := 1; tok.Line != want {
				t.Errorf("want line %v, got %v", want, tok.Line)
			}
			if tok.Kind != token.Call {
				t.Errorf("want kind %q, got %q", token.Call, tok.Kind)
			}
			if tok.Lexeme != tc.input {
				t.Errorf("want lexeme %q, got %q", tc.input, tok.Lexeme)
			}
		})
	}
}

func TestVariableTokens(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input string
	}{
		{"starts with underscore", "_foo"},
		{"starts with alpha", "foo"},
		{"contains underscore", "foo_bar"},
		{"contains digit", "foo123"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			tok, err := l.Consume()
			if err != nil {
				t.Fatal(err)
			}

			if want := 1; tok.Line != want {
				t.Errorf("want line %v, got %v", want, tok.Line)
			}
			if tok.Kind != token.Variable {
				t.Errorf("want kind %q, got %q", token.Variable, tok.Kind)
			}
			if tok.Lexeme != tc.input {
				t.Errorf("want lexeme %q, got %q", tc.input, tok.Lexeme)
			}
		})
	}
}

func TestLiteralTokens(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  string
		lexeme string
		kind   token.Kind
	}{
		{"integer simple", "1234567890", "1234567890", token.Integer},
		{"integer leading zeros", "0001234567890", "0001234567890", token.Integer},
		{"integer leading spaces", "   1234567890", "1234567890", token.Integer},
		{"integer trailing spaces", "1234567890    ", "1234567890", token.Integer},

		{"string double quotes", `"Hello, World!"`, "Hello, World!", token.String},
		{"string single quotes", `'Hello, World!'`, "Hello, World!", token.String},
		{"string double quotes empty", `''`, "", token.String},
		{"string single quotes empty", `""`, "", token.String},
		{"string with newline", "'Hello\nWorld!'", "Hello\nWorld!", token.String},

		// These are actually functions in the Knight spec, but we'll treat them as literals
		{"true (function)", "T", "T", token.True},
		{"true long (function)", "TRUE", "TRUE", token.True},
		{"false (function)", "F", "F", token.False},
		{"false long (function)", "FALSE", "FALSE", token.False},
		{"null (function)", "N", "N", token.Null},
		{"null long (function)", "NULL", "NULL", token.Null},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			tok, err := l.Consume()
			if err != nil {
				t.Fatal(err)
			}

			if want := 1; tok.Line != want {
				t.Errorf("want line %v, got %v", want, tok.Line)
			}
			if tok.Kind != tc.kind {
				t.Errorf("want kind %q, got %q", tc.kind, tok.Kind)
			}
			if tok.Lexeme != tc.lexeme {
				t.Errorf("want lexeme %q, got %q", tc.lexeme, tok.Lexeme)
			}
		})
	}
}

func TestUnknownTokens(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		input string
	}{
		{"symbol .", "."},
		{"unicode japanese", "エラー"},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			l.Load(strings.NewReader(tc.input))

			tok, err := l.Consume()
			if err == nil {
				t.Error("want error, got nil")
			}

			fmt.Println(tok, err)

			if tok.Kind != token.Unknown {
				t.Errorf("want kind %q, got %q", token.Unknown, tok.Kind)
			}
			if tok.Lexeme != tc.input {
				t.Errorf("want lexeme %q, got %q", tc.input, tok.Lexeme)
			}
		})
	}
}

func TestTokenStream(t *testing.T) {
	t.Parallel()

	const source = `
		# This comment should be ignored
		"Foo" # This comment should also be ignored
		B BLOCK C CALL D DUMP E EVAL F FALSE G GET I IF L LENGTH N NULL O OUTPUT P PROMPT Q QUIT R RAND S SUBSTITUTE
		T TRUE W WHILE X ! + - * / % ^ & | < > ? : ; =
		_foo foo foo_bar foo123
		1234567890 0001234567890 1234567890 1234567890
		"Hello, World!" 'Hello, World!' '' ""
'Hello
World!'
"Hello
World!"
		!!!foo +++bar ---baz ***qux ///abc %%%123 ^^^xxx &&&xxx |||xxx <<<xxx >>>xxx ???xxx ;;;xxx ===xxx
	`
	l := lexer.New()
	l.Load(strings.NewReader(source))

	var tokens []token.Token
	for tok, err := l.Consume(); tok.Kind != token.EOF; tok, err = l.Consume() {
		if err != nil {
			t.Fatal(err)
		}

		tokens = append(tokens, tok)
	}

	wants := []struct {
		lexeme string
		kind   token.Kind
	}{
		{"Foo", token.String},

		{"B", token.Call}, {"BLOCK", token.Call}, {"C", token.Call}, {"CALL", token.Call}, {"D", token.Call}, {"DUMP", token.Call},
		{"E", token.Call}, {"EVAL", token.Call}, {"F", token.False}, {"FALSE", token.False}, {"G", token.Call}, {"GET", token.Call},
		{"I", token.Call}, {"IF", token.Call}, {"L", token.Call}, {"LENGTH", token.Call}, {"N", token.Null}, {"NULL", token.Null},
		{"O", token.Call}, {"OUTPUT", token.Call}, {"P", token.Call}, {"PROMPT", token.Call}, {"Q", token.Call}, {"QUIT", token.Call},
		{"R", token.Call}, {"RAND", token.Call}, {"S", token.Call}, {"SUBSTITUTE", token.Call}, {"T", token.True}, {"TRUE", token.True},
		{"W", token.Call}, {"WHILE", token.Call}, {"X", token.Call}, {"!", token.Not}, {"+", token.Add}, {"-", token.Sub},
		{"*", token.Mul}, {"/", token.Div}, {"%", token.Mod}, {"^", token.Exp}, {"&", token.And}, {"|", token.Or}, {"<", token.Less},
		{">", token.Greater}, {"?", token.Equal} /*{":", token.Noop},*/, {";", token.Chain}, {"=", token.Assign},

		{"_foo", token.Variable}, {"foo", token.Variable}, {"foo_bar", token.Variable}, {"foo123", token.Variable},

		{"1234567890", token.Integer}, {"0001234567890", token.Integer}, {"1234567890", token.Integer}, {"1234567890", token.Integer},

		{"Hello, World!", token.String}, {"Hello, World!", token.String},
		{"", token.String}, {"", token.String},
		{"Hello\nWorld!", token.String}, {"Hello\nWorld!", token.String},

		{"!", token.Not}, {"!", token.Not}, {"!", token.Not}, {"foo", token.Variable},
		{"+", token.Add}, {"+", token.Add}, {"+", token.Add}, {"bar", token.Variable},
		{"-", token.Sub}, {"-", token.Sub}, {"-", token.Sub}, {"baz", token.Variable},
		{"*", token.Mul}, {"*", token.Mul}, {"*", token.Mul}, {"qux", token.Variable},
		{"/", token.Div}, {"/", token.Div}, {"/", token.Div}, {"abc", token.Variable},
		{"%", token.Mod}, {"%", token.Mod}, {"%", token.Mod}, {"123", token.Integer},
		{"^", token.Exp}, {"^", token.Exp}, {"^", token.Exp}, {"xxx", token.Variable},
		{"&", token.And}, {"&", token.And}, {"&", token.And}, {"xxx", token.Variable},
		{"|", token.Or}, {"|", token.Or}, {"|", token.Or}, {"xxx", token.Variable},
		{"<", token.Less}, {"<", token.Less}, {"<", token.Less}, {"xxx", token.Variable},
		{">", token.Greater}, {">", token.Greater}, {">", token.Greater}, {"xxx", token.Variable},
		{"?", token.Equal}, {"?", token.Equal}, {"?", token.Equal}, {"xxx", token.Variable},
		{";", token.Chain}, {";", token.Chain}, {";", token.Chain}, {"xxx", token.Variable},
		{"=", token.Assign}, {"=", token.Assign}, {"=", token.Assign}, {"xxx", token.Variable},

		// EOF token is ignored since if it didn't exist we'd infinitely loop anyway
	}
	if len(wants) != len(tokens) {
		t.Fatalf("want %v tokens, got %v", len(wants), len(tokens))
	}

	for i, want := range wants {
		if want.lexeme != tokens[i].Lexeme {
			t.Errorf("token [%v]: want lexeme %q, got %q", i, want.lexeme, tokens[i].Lexeme)
		}
		if want.kind != tokens[i].Kind {
			t.Errorf("token [%v] (%q): want kind %q, got %q", i, want.lexeme, want.kind, tokens[i].Kind)
		}
	}
}

func TestProgram(t *testing.T) {
	t.Parallel()

	b, err := os.ReadFile("testdata/knight.kn")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New()
	l.Load(bytes.NewReader(b))

	for tok, err := l.Consume(); tok.Kind != token.EOF; tok, err = l.Consume() {
		if err != nil {
			t.Fatal(err)
		}
	}
}
