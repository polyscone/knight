package value

import "fmt"

// Block is the result of the evaluation of the BLOCK function, and wraps an
// expression that is to be used as the argument to CALL.
type Block struct {
	Expr

	Value Expression
}

// AsBool is only implemented here so that Block can be used as a value.
// Any attempt to actually call this method will result in a panic, because a
// block should only be used for its expression as an argument to CALL.
func (b *Block) AsBool() *Bool {
	panic("a block cannot be converted into a bool")
}

// AsInt is only implemented here so that Block can be used as a value.
// Any attempt to actually call this method will result in a panic, because a
// block should only be used for its expression as an argument to CALL.
func (b *Block) AsInt() *Int {
	panic("a block cannot be converted into an int")
}

// AsString is only implemented here so that Block can be used as a value.
// Any attempt to actually call this method will result in a panic, because a
// block should only be used for its expression as an argument to CALL.
func (b *Block) AsString() *String {
	panic("a block cannot be converted into a string")
}

// AsExpr returns the wrapped Expression, which should be used as an argument
// to CALL.
func (b *Block) AsExpr() Expression {
	return b.Value
}

// Dump prints a string form of the Block for testing.
// A block doesn't actually need to dump anything according to the Knight spec
// but this implementation prints a representation of it anyway.
func (b *Block) Dump() string {
	if v, ok := b.Value.(Value); ok {
		return fmt.Sprintf("Block(%v)", v.Dump())
	}

	return "Block(AST)"
}

// String prints a string form of the Block as an s-expression for testing.
// The AsString method should be used to convert a value to a runtime String.
func (b *Block) String() string {
	return fmt.Sprintf("(block %v)", b.Value)
}

// NewBlock will return a Block that wraps the given expression, ready for use
// as an argument to CALL.
func NewBlock(v Expression) *Block {
	return &Block{Value: v}
}
