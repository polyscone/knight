package parser_test

import (
	"strings"
	"testing"

	"github.com/polyscone/knight/lexer"
	"github.com/polyscone/knight/parser"
	"github.com/polyscone/knight/value"
)

func TestGeneratedASTStructure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		source string
		want   string
	}{
		{
			"logical and (&)",
			"& 0 1",
			"(program (and 0 1))",
		},
		{
			"logical or (|)",
			"| 0 1",
			"(program (or 0 1))",
		},
		{
			"negation (!)",
			"! 0",
			"(program (not 0))",
		},
		{
			"add (+)",
			"+ 0 1",
			"(program (add 0 1))",
		},
		{
			"sub (-)",
			"- 0 1",
			"(program (sub 0 1))",
		},
		{
			"mul (*)",
			"* 0 1",
			"(program (mul 0 1))",
		},
		{
			"div (/)",
			"/ 0 1",
			"(program (div 0 1))",
		},
		{
			"modulo (%)",
			"% 0 1",
			"(program (mod 0 1))",
		},
		{
			"less than (<)",
			"< 0 1",
			"(program (less 0 1))",
		},
		{
			"greater than (>)",
			"> 0 1",
			"(program (greater 0 1))",
		},
		{
			"assign (=)",
			"= 0 1",
			"(program (assign 0 1))",
		},
		{
			"equal (?)",
			"? 0 1",
			"(program (equal 0 1))",
		},
		{
			"exponentiation (^)",
			"^ 0 1",
			"(program (exp 0 1))",
		},
		{
			"system (`)",
			"` \"foo bar\"",
			`(program (system "foo bar"))`,
		},
		{
			"chaining (;)",
			"; 0 1",
			"(program (chain 0 1))",
		},
		// {
		// 	"no-op (:)",
		// 	": 1",
		// 	"(program (noop 1))",
		// },
		{
			"word function block",
			"B 1",
			"(program (B 1))",
		},
		{
			"word function block long",
			"BLOCK 1",
			"(program (BLOCK 1))",
		},
		{
			"word function block long expression",
			"B + 1 2",
			"(program (B (add 1 2)))",
		},
		{
			"word function call",
			"C B 1",
			"(program (C (B 1)))",
		},
		{
			"word function call long",
			"CALL BLOCK 1",
			"(program (CALL (BLOCK 1)))",
		},
		{
			"word function dump",
			"D 1",
			"(program (D 1))",
		},
		{
			"word function dump long",
			"DUMP 1",
			"(program (DUMP 1))",
		},
		{
			"word function eval",
			`E "1"`,
			`(program (E "1"))`,
		},
		{
			"word function eval long",
			`EVAL "1"`,
			`(program (EVAL "1"))`,
		},
		{
			"word function get",
			`G 'foo' 0 1`,
			`(program (G "foo" 0 1))`,
		},
		{
			"word function get long",
			`GET 'foo' 0 1`,
			`(program (GET "foo" 0 1))`,
		},
		{
			"word function if",
			"I 0 1 2",
			"(program (I 0 1 2))",
		},
		{
			"word function if long",
			"IF 0 1 2",
			"(program (IF 0 1 2))",
		},
		{
			"word function length",
			`L "foo"`,
			`(program (L "foo"))`,
		},
		{
			"word function length long",
			`LENGTH "foo"`,
			`(program (LENGTH "foo"))`,
		},
		{
			"word function output",
			`O "foo"`,
			`(program (O "foo"))`,
		},
		{
			"word function output long",
			`OUTPUT "foo"`,
			`(program (OUTPUT "foo"))`,
		},
		{
			"word function prompt",
			"P",
			"(program P)",
		},
		{
			"word function prompt long",
			"PROMPT",
			"(program PROMPT)",
		},
		{
			"word function quit",
			"Q 1",
			"(program (Q 1))",
		},
		{
			"word function quit long",
			"QUIT 1",
			"(program (QUIT 1))",
		},
		{
			"word function rand",
			"R",
			"(program R)",
		},
		{
			"word function rand long",
			"RAND",
			"(program RAND)",
		},
		{
			"word function substitute",
			`S "foo" 0 1 "b"`,
			`(program (S "foo" 0 1 "b"))`,
		},
		{
			"word function substitute long",
			`SUBSTITUTE "foo" 0 1 "b"`,
			`(program (SUBSTITUTE "foo" 0 1 "b"))`,
		},
		{
			"word function set long",
			`SET "foo" 0 1 "b"`,
			`(program (SET "foo" 0 1 "b"))`,
		},
		{
			"word function while",
			"W 1 B 2",
			"(program (W 1 (B 2)))",
		},
		{
			"word function while long",
			"WHILE 1 BLOCK 2",
			"(program (WHILE 1 (BLOCK 2)))",
		},
		{
			"true (function)",
			"T",
			"(program true)",
		},
		{
			"true long (function)",
			"TRUE",
			"(program true)",
		},
		{
			"false (function)",
			"F",
			"(program false)",
		},
		{
			"false long (function)",
			"FALSE",
			"(program false)",
		},
		{
			"null (function)",
			"N",
			"(program null)",
		},
		{
			"null long (function)",
			"NULL",
			"(program null)",
		},
		{
			"_foo",
			"_foo",
			`(program (var "_foo"))`,
		},
		{
			"foo",
			"foo",
			`(program (var "foo"))`,
		},
		{
			"foo_bar",
			"foo_bar",
			`(program (var "foo_bar"))`,
		},
		{
			"foo123",
			"foo123",
			`(program (var "foo123"))`,
		},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New()
			p := parser.New(l)
			g := value.NewGlobalStore()
			program, err := p.Parse(g, strings.NewReader(tc.source))
			if err != nil {
				t.Fatal(err)
			}

			if program.String() != tc.want {
				t.Errorf("\nwant:\n\t%s\ngot:\n\t%s", tc.want, program)
			}
		})
	}
}
