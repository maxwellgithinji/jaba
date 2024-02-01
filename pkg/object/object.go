/*
* Package object helps represent the values encountered when evaluating the jaba program as an object.
* Every value will be wrapped in a struct that fulfills the object interface.
* The object system leverages on the host language (Go) data types and formatting methods to represent its values
 */
package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/maxwellgithinji/jaba/pkg/ast"
)

// ObjectType represents the category of the object
type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION_OBJECT"
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

// ReturnValue represents a jaba return value
// It fulfills the object interface by implementing the Type() and Inspect() methods
type ReturnValue struct {
	Value Object
}

// Type returns the type of the object
func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

// Inspect returns the string representation of the object value, return value
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

// Error represents internal jaba error
// it fulfills the Object interface by implementing the Type() and Inspect() methods
type Error struct {
	Message string
}

// Type returns the type of the object, error
func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

// Inspect returns the string representation of the object value, error
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Function represents a jaba function and may include parameters and some statements to be executed
// it fulfills the Object interface by implementing the Type() and Inspect() methods
type Function struct {
	// Parameters is a list of identifiers that should be passed to the function call
	Parameters []*ast.Identifier

	// Body contains a list of function statements to be evaluated
	Body *ast.BlockStatement

	// Env keeps track of variables during interpreter execution
	Env *Environment
}

// Type returns the type of the object, function
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

// Inspect returns the string representation of the function
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
