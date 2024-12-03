package evaluator

import (
	"github.com/kev/object"
)

// Builtins is a map of built-in functions
var builtins = map[string]*object.Builtin{

	// len function returns the length of the object
	// passed to it. It only supports strings... For now.
	"len": {
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
