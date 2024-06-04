package evaluator

import (
	"1ylang/object"
	"fmt"
	"math/rand"
)

// newBuiltin is a helper function to create a new builtin function object.
func newBuiltin(fn func(args ...object.Object) object.Object) *object.Builtin {
	return &object.Builtin{Fn: fn}
}

var builtins = map[string]*object.Builtin{
	"lang": newBuiltin(func(args ...object.Object) object.Object {
		WELCOME := `Welcome to 1lang! This is the 1lang programming language!`
		return &object.String{Value: WELCOME}
	}),
	"len": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	}),
	"puts": newBuiltin(func(args ...object.Object) object.Object {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return NULL
	}),
	"first": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}

		return NULL
	}),
	"last": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1]
		}

		return NULL
	}),
	"rest": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[1:length])
			return &object.Array{Elements: newElements}
		}

		return NULL
	}),
	"push": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		newElements := make([]object.Object, length+1)
		copy(newElements, arr.Elements)
		newElements[length] = args[1]

		return &object.Array{Elements: newElements}
	}),
	"concat": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("wrong number of arguments. got=%d, want=2+", len(args))
		}

		for _, arg := range args {
			if arg.Type() != object.ARRAY_OBJ {
				return newError("argument to `concat` must be ARRAY, got %s", arg.Type())
			}
		}

		newElements := []object.Object{}
		for _, arg := range args {
			newElements = append(newElements, arg.(*object.Array).Elements...)
		}

		return &object.Array{Elements: newElements}
	}),
	"range": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
			return newError("arguments to `range` must be INTEGER, got %s and %s", args[0].Type(), args[1].Type())
		}

		start := args[0].(*object.Integer).Value
		end := args[1].(*object.Integer).Value

		if start > end {
			return newError("start index cannot be greater than end index")
		}

		newElements := []object.Object{}
		for i := start; i < end; i++ {
			newElements = append(newElements, &object.Integer{Value: i})
		}

		return &object.Array{Elements: newElements}
	}),
	"random": newBuiltin(func(args ...object.Object) object.Object {
		if len(args) == 1 {
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `random` must be INTEGER, got %s", args[0].Type())
			}
			max := args[0].(*object.Integer).Value
			return &object.Integer{Value: rand.Int63n(max)}
		} else if len(args) == 2 {
			if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
				return newError("arguments to `random` must be INTEGER, got %s and %s", args[0].Type(), args[1].Type())
			}
			min := args[0].(*object.Integer).Value
			max := args[1].(*object.Integer).Value
			return &object.Integer{Value: rand.Int63n(max-min) + min}
		} else {
			return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
		}
	}),
}
