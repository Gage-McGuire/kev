package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Gage-McGuire/kev/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
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

// Represents a string object
type String struct {
	Value string
}

// Represents a boolean object
type Boolean struct {
	Value bool
}

// Represents a null object
type Null struct{}

// Represents a return value object
type ReturnValue struct {
	Value Object
}

// Represents an error object
type Error struct {
	Message string
}

type Array struct {
	Elements []Object
}

// Represents a function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Hashable interface {
	HashKey() HashKey
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Returns the string representation
// of the function object
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Returns the type of the function object
// which is always a FUNCTION_OBJ
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

// Returns the value of the integer object
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Returns the type of the integer object
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Returns the value of the string object
func (s *String) Inspect() string {
	return s.Value
}

// Returns the type of the string object
func (s *String) Type() ObjectType {
	return STRING_OBJ
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

// Returns the value of the return value object
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Returns the type of the return value object
// which is always a RETURN_VALUE_OBJ
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// Returns the value of the error object
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Returns the type of the error object
// which is always a ERROR_OBJ
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// Returns the value of the array object
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Returns the type of the array object
func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

/*
 * Environment
 */

// Creates a new environment
// with an empty store
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Creates a new environment
// which is enclosed and limited to its block statement
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Represents an environment with a store
// that holds the bindings of the variables
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Returns the object with the given name
// contained in the store
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Sets the object with the given name
// in the store
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

/*
 * Built-in functions
 */

// Represents a built-in function
type BuiltinFunction func(args ...Object) Object

// Represents a built-in function object
type Builtin struct {
	Func BuiltinFunction
}

// Returns the value of the built-in function object
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Returns the Type of the built-in function object
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
