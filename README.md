![Golang workflow](https://github.com/maxwellgithinji/jaba/actions/workflows/go.yml/badge.svg) 
[![wakatime](https://wakatime.com/badge/user/5af887ac-99ff-4b74-9e6a-34c9b421a9d6/project/018d23bd-78d1-4d92-baae-d1019d0fef51.svg)](https://wakatime.com/badge/user/5af887ac-99ff-4b74-9e6a-34c9b421a9d6/project/018d23bd-78d1-4d92-baae-d1019d0fef51)

# jaba
jaba is a programming language built using golang. The language source code is transformed into tokens through lexical analysis and then tokens are transformed to an Abstract Syntax Tree using a parser. Evaluation (Interpretation) is done to give the AST meaning and returns the respective result of the source code entered.

Source Code -> Tokens -> Abstract Syntax tree (AST) -> Output

The first transformation, from source code to tokens, is called “lexical analysis”, or “lexing” for short. It’s done by a lexer (also called tokenizer or scanner – some use one word or the other to denote subtle differences in behaviour). 
Tokens themselves are small, easily categorizable data structures that are then fed to the parser, which does the second transformation and turns the tokens into an “Abstract Syntax Tree”. The final step is evaluation where the AST is given meaning by an interpreter which gives the respective result. The interpreter used here is called a `tree walking interpreter`. An object is used for the interpreter; it is not very performant, but it's easy to get started with.

## Supported Features
- C-like syntax
- variable bindings
- integers and booleans
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures

## Getting Started

### installing
1. clone the repository
2. run `go run main.go`
3. Enter the jaba program on the command line


## Examples 

### Variable binding
```
let age = 1;
let name = "Trent";
let result = 10 * (20 / 2);
```
### Accessing Elements
```
let myArray = [1, 2, 3, 4, 5];
myArray[0]       // => 1

let thorsten = {"name": "Thorsten", "age": 28};
thorsten["name"] // => "Thorsten"
```
### Function Binding
```
let add = fn(a, b) { return a + b; };
```

### Implicit Returns
```
let add = fn(a, b) { a + b; };
```

### Function Call
```
add(1, 2);
```
### Complex Function
```
let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      1
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
```

### Higher Order Functions
```
let twice = fn(f, x) { return f(f(x)); };

let addTwo = fn(x) { return x + 2; };

twice(addTwo, 2); // => 6
```

### Closures
```
// newGreeter returns a new function, that greets a `name` with the given
// `greeting`.
let newGreeter = fn(greeting) {
  // `puts` is a built-in function we add to the interpreter
  return fn(name) { puts(greeting + " " + name); }
};

// `hello` is a greeter function that says "Hello"
let hello = newGreeter("Hello");

// Calling it outputs the greeting:
hello("dear, future Reader!"); // => Hello dear, future Reader!
```

## Garbage Collection
The host language for jaba is Golang, which does the garbage collection and allows the memory not to leak when we run code like this
```
let counter = fn(x) {
  if (x > 100) {
    return true;
  } else {
    let foobar = 9999;
    counter(x + 1); 
  }
}; 

counter(0);
```
If C was used to write the interpreter, we would have to implement our own garbage collector to avoid memory leaks. 
An example of a garbage collector implementation is [Mark and Sweep](https://www.geeksforgeeks.org/mark-and-sweep-garbage-collection-algorithm/).

## Upcoming feature
1. Evolve the REPL to launch as a cmd application
2. Support reading of `.jaba` files for good programming execution
2. support syntax highlighting
4. Language documentation

## References 
### [ Writing An Interpreter In Go](https://interpreterbook.com/)
- Ball, Thorsten. Writing An Interpreter In Go (p. 8). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (pp. 8-9). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (p. 9). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (p. 13). Kindle Edition. 
- Ball, Thorsten. Writing An Interpreter In Go (pp. 189-190). Kindle Edition. 
