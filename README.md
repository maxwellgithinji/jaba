# jaba
jaba is a programming language built using golang. The language source code is transformed into tokens through lexical analysis and then tokens are transformed to an Abstract Syntax Tree using a parser

Source Code -> Tokens -> Abstract Syntax tree (AST)

The first transformation, from source code to tokens, is called “lexical analysis”, or “lexing” for short. It’s done by a lexer (also called tokenizer or scanner – some use one word or the other to denote subtle differences in behaviour). 
Tokens themselves are small, easily categorizable data structures that are then fed to the parser, which does the second transformation and turns the tokens into an “Abstract Syntax Tree”.

## Supported Features
- C-like syntax
- variable bindings
- integers and booleans
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures

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

## References 
### [ Writing An Interpreter In Go](https://interpreterbook.com/)
- Ball, Thorsten. Writing An Interpreter In Go (p. 8). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (pp. 8-9). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (p. 9). Kindle Edition.
- Ball, Thorsten. Writing An Interpreter In Go (p. 13). Kindle Edition. 
