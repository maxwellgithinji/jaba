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
	STRING_OBJECT       = "STRING"
	BUILTIN_OBJECT      = "BUILTIN"
	ARRAY_OBJECT        = "ARRAY"
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

// String represents a jaba string which is an expression which evaluates to a value
// it fulfills the Object interface by implementing the Type() and Inspect() methods
type String struct {
	// Value is the actual value of the string literal
	Value string
}

// Type returns the type of the object, string
func (s *String) Type() ObjectType {
	return STRING_OBJECT
}

// Inspect returns the string representation of the object value, string
func (s *String) Inspect() string {
	return s.Value
}

// BuiltinFunction represents a jaba builtin function which is from the host language that allows user to
// use host language functions
type BuiltinFunction func(args ...Object) Object

// Builtin is a wrapper around golang function which is the host language
// it fulfills the Object interface by implementing the Type() and Inspect() methods
type Builtin struct {
	Function BuiltinFunction
}

// Type returns the type of the object, builtin
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJECT
}

// Inspect returns the string representation of the object value, builtin
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Array represents a jaba builtin array of objects
// it fulfills the Object interface by implementing the Type() and Inspect() methods
type Array struct {
	Elements []Object
}

// Type returns the type of the object, array
func (a *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

// Inspect returns the string representation of the object value, array
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
