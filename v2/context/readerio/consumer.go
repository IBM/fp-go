package readerio

import "github.com/IBM/fp-go/v2/io"

// ChainConsumer chains a consumer function into a ReaderIO computation, discarding the original value.
// This is useful for performing side effects (like logging or metrics) that consume a value
// but don't produce a meaningful result.
//
// The consumer is executed for its side effects, and the computation returns an empty struct.
//
// Type Parameters:
//   - A: The type of value to consume
//
// Parameters:
//   - c: A consumer function that performs side effects on the value
//
// Returns:
//   - An Operator that chains the consumer and returns struct{}
//
// Example:
//
//	logUser := func(u User) {
//	    log.Printf("Processing user: %s", u.Name)
//	}
//
//	pipeline := F.Pipe2(
//	    fetchUser(123),
//	    ChainConsumer(logUser),
//	)
//
//go:inline
func ChainConsumer[A any](c Consumer[A]) Operator[A, Void] {
	return ChainIOK(io.FromConsumer(c))
}

// ChainFirstConsumer chains a consumer function into a ReaderIO computation, preserving the original value.
// This is useful for performing side effects (like logging or metrics) while passing the value through unchanged.
//
// The consumer is executed for its side effects, but the original value is returned.
//
// Type Parameters:
//   - A: The type of value to consume and return
//
// Parameters:
//   - c: A consumer function that performs side effects on the value
//
// Returns:
//   - An Operator that chains the consumer and returns the original value
//
// Example:
//
//	logUser := func(u User) {
//	    log.Printf("User: %s", u.Name)
//	}
//
//	pipeline := F.Pipe3(
//	    fetchUser(123),
//	    ChainFirstConsumer(logUser),  // Logs but passes user through
//	    Map(func(u User) string { return u.Email }),
//	)
//
//go:inline
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A] {
	return ChainFirstIOK(io.FromConsumer(c))
}
