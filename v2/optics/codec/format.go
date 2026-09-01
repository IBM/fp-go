package codec

import (
	"fmt"
	"log/slog"
)

// String implements the fmt.Stringer interface for Type.
// Returns the name of the type for debugging and display purposes.
// The exact format is not a stable contract and may change across versions.
//
// Example:
//
//	stringType := codec.String()
//	fmt.Println(stringType) // Output: "string"
func (t Type[A, O, I]) String() string {
	return t.name
}

// Format implements the fmt.Formatter interface for Type.
// It provides custom formatting based on the format verb:
//   - %s, %v: Returns the type name
//   - %q: Returns the type name in quotes
//   - other verbs: Returns the type name
//
// The exact output format is not a stable contract and may change across versions.
func (t Type[A, O, I]) Format(f fmt.State, verb rune) {
	switch verb {
	case 'q':
		fmt.Fprintf(f, "%q", t.name)
	default:
		fmt.Fprint(f, t.name)
	}
}

// LogValue implements the slog.LogValuer interface for Type.
// Returns a slog.Value containing the type name for structured logging.
// The exact structure of the returned slog.Value is not a stable contract and may change across versions.
//
// Example:
//
//	stringType := codec.String()
//	slog.Info("codec created", "codec", stringType)
//	// Logs: codec={name=string}
func (t Type[A, O, I]) LogValue() slog.Value {
	return slog.StringValue(t.name)
}
