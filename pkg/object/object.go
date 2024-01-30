/*
* Package object helps represent the values encountered when evaluating the jaba program as an object.
* Every value will be wrapped in a struct that fulfills the object interface.
* The object system leverages on the host language (Go) data types and formatting methods to represent its values
 */
package object

import "fmt"

// ObjectType represents the category of the object
type ObjectType string

const (
	INTEGER_OBJECT = "INTEGER"
	BOOLEAN_OBJECT = "BOOLEAN"
	NULL_OBJECT    = "NULL"
)

// Object is an interface that helps represent the values encountered when evaluating the jaba program
type Object interface {
	// Type returns the type of the object
	Type() ObjectType

	// Inspect returns the string representation of the object value
	Inspect() string
}

// Integer is a jaba data type that represents numbers
// It fulfills the object interface by implementing the Type() and Inspect() methods
type Integer struct {
	Value int64
}

// Type returns the type of the object
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

// Inspect returns the string representation of the object value, integer
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean is a jaba data type that represents true or false
// It fulfills the object interface by implementing the Type() and Inspect() methods
type Boolean struct {
	Value bool
}

// Type returns the type of the object
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

// Inspect returns the string representation of the object value, boolean
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null represents absence of a value
// It fulfills the object interface by implementing the Type() and Inspect() methods
type Null struct {
	Value interface{}
}

// Type returns the type of the object
func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

// Inspect returns the string representation of the object value, null
func (n *Null) Inspect() string {
	return "null"
}