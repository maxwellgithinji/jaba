/*
* Package object helps represent the values encountered when evaluating the jaba program as an object.
* Every value will be wrapped in a struct that fulfills the object interface.
* The object system leverages on the host language (Go) data types and formatting methods to represent its values
 */
package object

// Environment is a wrapper of the map implementation that helps associate a string key with an object
type Environment struct {
	// store is the hashmap that stores the objects
	store map[string]Object

	// outer helps with scoping of the environment.
	// its helpful when separating program and function variables
	outer *Environment
}

// NewEnvironment creates a new instance of the environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new instance of an scoped environment
// it extends the Environment instance to cater for enclosed environments
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()

	env.outer = outer

	return env
}

// Get returns the object associated with the given key from the environment
// it also checks for values both in the inner and outer scopes
func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}

	return obj, ok
}

// Set creates an object in the environment hashmap and returns what was created
func (e *Environment) Set(key string, value Object) Object {
	e.store[key] = value
	return value
}
