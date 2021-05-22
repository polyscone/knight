package ast

import "github.com/polyscone/knight/value"

func NewGlobal(v *value.Global) Node {
	return v
}

func NewBool(v bool) Node {
	return value.NewBool(v)
}

func NewInt(v int) Node {
	return value.NewInt(v)
}

func NewString(v string) Node {
	return value.NewString(v)
}

func NewNull() Node {
	return value.NewNull()
}
