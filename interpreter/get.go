package interpreter

import "github.com/polyscone/knight/value"

// Get returns a substring of the given string value.
func (i *Interpreter) Get(strVal, startVal, countVal value.Value) (value.Value, error) {
	start := startVal.AsInt().Value
	count := countVal.AsInt().Value
	amount := start + count

	if amount == 0 {
		return value.NewString(""), nil
	}

	str := strVal.AsString()

	if start == 0 && amount == len(str.Value) {
		return strVal, nil
	}

	return value.NewString(str.Value[start:amount]), nil
}
