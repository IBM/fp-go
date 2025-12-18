package formatting

import (
	"fmt"
	"log/slog"
)

type (

	// Formattable is a composite interface that combines multiple formatting capabilities
	// from the Go standard library. Types implementing this interface can be formatted
	// in various contexts including string conversion, custom formatting, Go syntax
	// representation, and structured logging.
	//
	// This interface is particularly useful for types that need to provide consistent
	// formatting across different output contexts, such as logging, debugging, and
	// user-facing displays.
	//
	// Embedded Interfaces:
	//
	//   - fmt.Stringer: Provides String() string method for basic string representation
	//   - fmt.Formatter: Provides Format(f fmt.State, verb rune) for custom formatting with verbs like %v, %s, %+v, etc.
	//   - fmt.GoStringer: Provides GoString() string method for Go-syntax representation (used with %#v)
	//   - slog.LogValuer: Provides LogValue() slog.Value for structured logging with the slog package
	//
	// Example Implementation:
	//
	//	type User struct {
	//		ID   int
	//		Name string
	//	}
	//
	//	// String provides a simple string representation
	//	func (u User) String() string {
	//		return fmt.Sprintf("User(%s)", u.Name)
	//	}
	//
	//	// Format provides custom formatting based on the verb
	//	func (u User) Format(f fmt.State, verb rune) {
	//		switch verb {
	//		case 'v':
	//			if f.Flag('+') {
	//				fmt.Fprintf(f, "User{ID: %d, Name: %s}", u.ID, u.Name)
	//			} else {
	//				fmt.Fprint(f, u.String())
	//			}
	//		case 's':
	//			fmt.Fprint(f, u.String())
	//		}
	//	}
	//
	//	// GoString provides Go-syntax representation
	//	func (u User) GoString() string {
	//		return fmt.Sprintf("User{ID: %d, Name: %q}", u.ID, u.Name)
	//	}
	//
	//	// LogValue provides structured logging representation
	//	func (u User) LogValue() slog.Value {
	//		return slog.GroupValue(
	//			slog.Int("id", u.ID),
	//			slog.String("name", u.Name),
	//		)
	//	}
	//
	// Usage:
	//
	//	user := User{ID: 1, Name: "Alice"}
	//	fmt.Println(user)           // Output: User(Alice)
	//	fmt.Printf("%+v\n", user)   // Output: User{ID: 1, Name: Alice}
	//	fmt.Printf("%#v\n", user)   // Output: User{ID: 1, Name: "Alice"}
	//	slog.Info("user", "user", user) // Structured log with id and name fields
	Formattable interface {
		fmt.Stringer
		fmt.Formatter
		fmt.GoStringer
		slog.LogValuer
	}
)
