package codec

import (
	"fmt"
	"log/slog"
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
//   - other verbs: Returns the type name
func (t *typeImpl[A, O, I]) Format(f fmt.State, verb rune) {
	switch verb {
	case 'q':
		fmt.Fprintf(f, "%q", t.name)
	default:
		fmt.Fprint(f, t.name)
	}
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
//	// Logs: codec={name=string}
func (t *typeImpl[A, O, I]) LogValue() slog.Value {
	return slog.StringValue(t.name)
}
