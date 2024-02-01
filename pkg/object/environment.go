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
}

// NewEnvironment creates a new instance of the environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s}
}

// Get returns the object associated with the given key from the environment
func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	return obj, ok
}

// Set creates an object in the environment hashmap and returns what was created
func (e *Environment) Set(key string, value Object) Object {
	e.store[key] = value
	return value
}
