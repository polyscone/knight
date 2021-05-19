package interpreter

import (
	"strings"

	"github.com/polyscone/knight/value"
)

type StringReader interface {
	ReadString(delim byte) (string, error)
}

// Prompt will ask the user for some input on stdin and return a string value.
func (i *Interpreter) Prompt(r StringReader) (value.Value, error) {
	text, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return value.NewString(strings.TrimRight(text, "\r\n")), nil
}
