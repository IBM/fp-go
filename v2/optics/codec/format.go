package codec

import (
	"fmt"
	"log/slog"

	"github.com/IBM/fp-go/v2/internal/formatting"
)

// String implements the fmt.Stringer interface for typeImpl.
// It returns the name of the type, which is used for simple string representation.
//
// Example:
//
//	stringType := codec.String()
//	fmt.Println(stringType) // Output: "string"
func (t *typeImpl[A, O, I]) String() string {
	return t.name
}

// Format implements the fmt.Formatter interface for typeImpl.
// It provides custom formatting based on the format verb:
//   - %s, %v: Returns the type name
//   - %q: Returns the type name in quotes
//   - %#v: Returns a detailed Go-syntax representation
//
// Example:
//
//	intType := codec.Int()
//	fmt.Printf("%s\n", intType)   // Output: int
//	fmt.Printf("%q\n", intType)   // Output: "int"
//	fmt.Printf("%#v\n", intType)  // Output: codec.Type[int, int, any]{name: "int"}
func (t *typeImpl[A, O, I]) Format(f fmt.State, verb rune) {
	formatting.FmtString(t, f, verb)
}

// GoString implements the fmt.GoStringer interface for typeImpl.
// It returns a Go-syntax representation of the type that could be used
// to recreate the type (though not executable due to function values).
//
// This is called when using the %#v format verb with fmt.Printf.
//
// Example:
//
//	stringType := codec.String()
//	fmt.Printf("%#v\n", stringType)
//	// Output: codec.Type[string, string, any]{name: "string"}
func (t *typeImpl[A, O, I]) GoString() string {
	return fmt.Sprintf("codec.Type[%s, %s, %s]{name: %q}",
		typeNameOf[A](), typeNameOf[O](), typeNameOf[I](), t.name)
}

// LogValue implements the slog.LogValuer interface for typeImpl.
// It provides structured logging representation of the codec type.
// Returns a slog.Value containing the type information as a group with
// the codec name and type parameters.
//
// This method is called automatically when logging a codec with slog.
//
// Example:
//
//	stringType := codec.String()
//	slog.Info("codec created", "codec", stringType)
//	// Logs: codec={name=string type_a=string type_o=string type_i=interface {}}
func (t *typeImpl[A, O, I]) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", t.name),
		slog.String("type_a", typeNameOf[A]()),
		slog.String("type_o", typeNameOf[O]()),
		slog.String("type_i", typeNameOf[I]()),
	)
}

// typeNameOf returns a string representation of the type T.
// It handles the special case where T is 'any' (interface{}).
func typeNameOf[T any]() string {
	var zero T
	typeName := fmt.Sprintf("%T", zero)
	// Handle the case where %T prints "<nil>" for interface{} types
	if typeName == "<nil>" {
		return "interface {}"
	}
	return typeName
}
