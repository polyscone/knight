package ast

import (
	"strings"
	"unicode/utf8"

	"github.com/polyscone/knight/value"
)

// Invalid represents an invalid AST node that, if part of a larger AST, should
// render the entire tree invalid.
var Invalid Node

// Node is a wrapper around the lower level AST node.
// It also wraps concrete values in the case of leaf nodes.
type Node interface {
	value.Expression

	ASTString(string) string
}

func indent(str string, length int) string {
	lines := strings.Split(str, "\n")
	indented := make([]string, len(lines))

	for i, line := range lines {
		indented[i] = strings.Repeat(" ", length) + line
	}

	return strings.Join(indented, "\n")
}

func prefix(str, first, rest string) string {
	lines := strings.Split(str, "\n")
	prefixed := make([]string, len(lines))

	if rest == "" {
		rest = strings.Repeat(" ", utf8.RuneCountInString(first))
	}

	for i, line := range lines {
		if i == 0 {
			prefixed[i] = first + line
		} else {
			prefixed[i] = rest + line
		}
	}

	return strings.Join(prefixed, "\n")
}

func SprintNode(style, name string, parts ...string) string {
	switch style {
	case "waterfall":
		pretty := name

		if len(parts) > 0 {
			pretty += " ┐"
			length := utf8.RuneCountInString(pretty) - 1

			for i, str := range parts {
				if isLast := i == len(parts)-1; isLast {
					pretty += "\n" + indent(prefix(str, "└ ", ""), length)
				} else {
					pretty += "\n" + indent(prefix(str, "├ ", "│ "), length)
				}
			}
		}

		return strings.TrimSpace(pretty)
	case "tree":
		pretty := name

		if len(parts) > 0 {
			for i, str := range parts {
				if isLast := i == len(parts)-1; isLast {
					pretty += "\n" + indent(prefix(str, "└─ ", ""), 0)
				} else {
					pretty += "\n" + indent(prefix(str, "├─ ", "│  "), 0)
				}
			}
		}

		return strings.TrimSpace(pretty)
	default:
		pretty := name

		if len(parts) > 0 {
			for _, str := range parts {
				pretty += " " + str
			}
		}

		return "(" + strings.TrimSpace(pretty) + ")"
	}
}
