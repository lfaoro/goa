package dsl

import (
	"goa.design/goa/design"
	"goa.design/goa/eval"
)

// Payload defines the data type of an method input. Payload also makes the
// input required.
//
// Payload must appear in a Method expression.
//
// Payload takes one to three arguments. The first argument is either a type or
// a DSL function. If the first argument is a type then an optional description
// may be passed as second argument. Finally a DSL may be passed as last
// argument that further specializes the type by providing additional
// validations (e.g. list of required attributes)
//
// The valid usage for Payload are thus:
//
//    Payload(Type)
//
//    Payload(func())
//
//    Payload(Type, "description")
//
//    Payload(Type, func())
//
//    Payload(Type, "description", func())
//
// Examples:
//
//    Method("upper"), func() {
//        // Use primitive type.
//        Payload(String)
//    }
//
//    Method("upper"), func() {
//        // Use primitive type.and description
//        Payload(String, "string to convert to uppercase")
//    }
//
//    Method("upper"), func() {
//        // Use primitive type, description and validations
//        Payload(String, "string to convert to uppercase", func() {
//            Pattern("^[a-z]")
//        })
//    }
//
//    Method("add", func() {
//        // Define payload data structure inline
//        Payload(func() {
//            Description("Left and right operands to add")
//            Attribute("left", Int32, "Left operand")
//            Attribute("right", Int32, "Left operand")
//            Required("left", "right")
//        })
//    })
//
//    Method("add", func() {
//        // Define payload type by reference to user type
//        Payload(Operands)
//    })
//
//    Method("divide", func() {
//        // Specify additional required attributes on user type.
//        Payload(Operands, func() {
//            Required("left", "right")
//        })
//    })
//
func Payload(val interface{}, args ...interface{}) {
	if len(args) > 2 {
		eval.ReportError("too many arguments")
	}
	e, ok := eval.Current().(*design.MethodExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	e.Payload = methodDSL("Payload", val, args...)
}

func methodDSL(suffix string, p interface{}, args ...interface{}) *design.AttributeExpr {
	var (
		att *design.AttributeExpr
		fn  func()
	)
	switch actual := p.(type) {
	case func():
		fn = actual
		att = &design.AttributeExpr{Type: &design.Object{}}
	case design.UserType:
		if len(args) == 0 {
			// Do not duplicate type if it is not customized
			return &design.AttributeExpr{Type: actual}
		}
		att = &design.AttributeExpr{Type: design.Dup(actual)}
	case design.DataType:
		att = &design.AttributeExpr{Type: actual}
	default:
		eval.ReportError("invalid %s argument, must be a type or a function", suffix)
		return nil
	}
	if len(args) >= 1 {
		if f, ok := args[len(args)-1].(func()); ok {
			if fn != nil {
				eval.ReportError("invalid arguments in %s call, must be (type), (func), (type, func), (type, desc) or (type, desc, func)", suffix)
			}
			fn = f
		}
		if d, ok := args[0].(string); ok {
			att.Description = d
		}
	}
	if fn != nil {
		eval.Execute(fn, att)
		if obj, ok := att.Type.(*design.Object); ok {
			if len(*obj) == 0 {
				att.Type = design.Empty
			}
		}
	}
	return att
}
