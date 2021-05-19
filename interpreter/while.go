package interpreter

import "github.com/polyscone/knight/value"

// While evaluates the given body expression for as long as the evaluation of
// the condition expression results in a truthy value.
func (i *Interpreter) While(condition, body value.Expression) (value.Value, error) {
	for {
		result, err := i.eval(condition)
		if err != nil {
			return nil, err
		}

		if !result.AsBool().Value {
			break
		}

		if _, err := i.eval(body); err != nil {
			return nil, err
		}
	}

	return value.NewNull(), nil
}
