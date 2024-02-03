/*
* Package evaluator uses the object system to evaluate the AST
 */
package evaluator

import (
	"github.com/maxwellgithinji/jaba/pkg/object"
)

// builtins is a hashmap to keep track of the variables during program execution
var builtins = map[string]*object.Builtin{
	"len": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d want: %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to len not supported, got: %s", args[0].Type())

			}
		},
	},
}
