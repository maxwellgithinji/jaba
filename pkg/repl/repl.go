/*
* Package repl (Read Eval Print Loop) or console is used to "Read" the input,
* sends it to the interpreter for "Evaluation", "Prints" the output of the interpreter, and then repeats the process("Loop").
 */
package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/maxwellgithinji/jaba/pkg/lexer"
	"github.com/maxwellgithinji/jaba/pkg/token"
)

// Prompt indicates the user start typing jaba code.
const Prompt = ">>"

// Run is a Read Eval Print Loop function that runs the jaba program.
// it helps the user code the jaba program on the command line
func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
