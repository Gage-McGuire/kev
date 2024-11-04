package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

// Base representation of an object.
// It holds the type of the object,
// and a string representation of the value of the object
type Object interface {
	// type of the object
	Type() ObjectType

	// value of the object
	Inspect() string
}

// Represents a integer object
type Integer struct {
	Value int64
}

// Represents a boolean object
type Boolean struct {
	Value bool
}

// Represents a null object
type Null struct{}

// Returns the value of the integer object
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Returns the type of the integer object
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Returns the value of the boolean object
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Returns the type of the boolean object
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// Returns the value of the null object
func (n *Null) Inspect() string {
	return "null"
}

// Returns the type of the null object
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}