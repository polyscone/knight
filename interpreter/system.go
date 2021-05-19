package interpreter

import (
	"os/exec"
	"runtime"

	"github.com/polyscone/knight/value"
)

// System executes the string conversion of the given value.
//
// Windows systems will use cmd, otherwise sh will be used instead.
func (i *Interpreter) System(val value.Value) (value.Value, error) {
	var sh string
	var args []string
	if runtime.GOOS == "windows" {
		sh = "cmd"
		args = append(args, "/c")
	} else {
		sh = "/bin/sh"
		args = append(args, "-c")
	}

	args = append(args, val.AsString().Value)

	out, err := exec.Command(sh, args...).Output()
	if err != nil {
		return nil, err
	}

	return value.NewString(string(out)), nil
}
