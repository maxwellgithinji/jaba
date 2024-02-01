/*
* Package repl (Read Eval Print Loop) or console is used to "Read" the input,
* sends it to the interpreter for "Evaluation", "Prints" the output of the interpreter, and then repeats the process("Loop").
 */
package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/maxwellgithinji/jaba/pkg/evaluator"
	"github.com/maxwellgithinji/jaba/pkg/lexer"
	"github.com/maxwellgithinji/jaba/pkg/object"
	"github.com/maxwellgithinji/jaba/pkg/parser"
)

// Prompt indicates the user start typing jaba code.
const Prompt = ">>"

// PRETTY_JABA a pretty printer that prints jaba logo
const PRETTY_JABA = `
____    
/oo  \   
|   __/    
/    _ |    
|     \ \    
\___  \ \__ 
|     \___\
`

// Run is a Read Eval Print Loop function that runs the jaba program.
// it helps the user code the jaba program on the command line
func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(out, Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, PRETTY_JABA)
	io.WriteString(out, "Woops! We ran into some jaba stories here!\n")
	io.WriteString(out, "parser errors: \n")
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
