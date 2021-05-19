package interpreter

import "github.com/polyscone/knight/value"

// Substitute returns a new string value where the given range is replaced with
// the given replacement value.
func (i *Interpreter) Substitute(strVal, startVal, countVal, replacementVal value.Value) (value.Value, error) {
	start := startVal.AsInt().Value
	count := countVal.AsInt().Value
	replacement := replacementVal.AsString()

	if count == 0 && replacement.Value == "" {
		return strVal, nil
	}

	str := strVal.AsString()
	amount := start + count

	if replacement.Value == "" {
		if start == 0 {
			if amount == len(str.Value) {
				return value.NewString(""), nil
			}

			return value.NewSubString(str, amount, len(str.Value)), nil
		}

		if amount == 0 || start == amount {
			return str, nil
		}

		if amount == len(str.Value) {
			return value.NewSubString(str, 0, start), nil
		}

		lhs := value.NewSubString(str, 0, start)
		rhs := value.NewSubString(str, amount, len(str.Value))

		return value.NewConcatString(lhs, rhs), nil
	}

	lhs := value.NewSubString(str, 0, start)
	lhs = value.NewConcatString(lhs, replacement)
	rhs := value.NewSubString(str, amount, len(str.Value))

	return value.NewConcatString(lhs, rhs), nil
}
